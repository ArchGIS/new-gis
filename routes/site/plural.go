package site

import (
	"net/http"
	"strings"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	item struct {
		Name     []string `json:"site_name"`
		ResName  []string `json:"research_name"`
		Epoch    int      `json:"epoch"`
		SiteType int      `json:"type"`
	}

	site struct {
		ID     uint64      `json:"id"`
		Item   item        `json:"item"`
		Coords coordinates `json:"coordinates"`
	}

	coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	requestParams struct {
		Name   string `query:"site_name"`
		Epoch  int    `query:"epoch_id" validate:"min=0,max=8"`
		Type   int    `query:"type_id" validate:"min=0,max=12"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}
)

const (
	statement = `
    MATCH (s:Monument)<--(k:Knowledge)
    MATCH (s)-[:has]->(st:MonumentType)
    MATCH (s)-[:has]->(e:Epoch)
    MATCH (r:Research)-[:has]->(k)
		%s
		WITH
			s.id as id,
			{
				site_name: collect(k.monument_name),
				research_name: collect(r.name),
				epoch: e.id,
				type: st.id
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`

	entity = "Monument"
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	sites, err := querySites(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"sites": sites})
}

func querySites(c echo.Context) (sites []site, err error) {
	req := &requestParams{
		Name:   "",
		Epoch:  0,
		Type:   0,
		Offset: 0,
		Limit:  20,
	}

	if err = c.Bind(req); err != nil {
		return nil, routes.NotAllowedQueryParams
	}

	if err = c.Validate(req); err != nil {
		return nil, routes.NotValidQueryParameters
	}

	cq := neo.BuildCypherQuery(
		cypher.Filter(statement, siteFilterString(req)),
		&sites,
		neoism.Props{
			"name":   cypher.BuildRegexpFilter(req.Name),
			"epoch":  req.Epoch,
			"type":   req.Type,
			"offset": req.Offset,
			"limit":  req.Limit,
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	if len(sites) > 0 {
		ids := make([]uint64, len(sites))
		coords := make([]coordinates, len(sites))

		for i, v := range sites {
			ids[i] = v.ID
		}

		coordsCQ := neo.BuildCypherQuery(
			cypher.BuildCoordinates(ids, entity, false),
			&coords,
			neoism.Props{},
		)

		err = neo.DB.Cypher(&coordsCQ)
		if err != nil {
			return nil, err
		}

		for i := range coords {
			sites[i].Coords = coords[i]
		}
	}

	return sites, nil
}

func siteFilterString(reqParams *requestParams) string {
	var filter []string
	var stmt string

	if reqParams.Name != "" {
		filter = append(filter, "k.monument_name =~ {name}")
	}
	if reqParams.Epoch != 0 {
		filter = append(filter, "e.id = {epoch}")
	}
	if reqParams.Type != 0 {
		filter = append(filter, "st.id = {type}")
	}
	if len(filter) > 0 {
		stmt = "WHERE " + strings.Join(filter, " AND ")
	}

	return stmt
}
