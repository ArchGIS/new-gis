package sites

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ArchGIS/new-gis/assert"
	cypher "github.com/ArchGIS/new-gis/cypherBuilder"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	item struct {
		Name     []string `json:"site_name"`
		ResName  []string `json:"research_name"`
		Epoch    string   `json:"epoch"`
		SiteType string   `json:"type"`
	}

	site struct {
		ID     int         `json:"id"`
		Item   item        `json:"item"`
		Coords coordinates `json:"coordinates"`
	}

	coordinates struct {
		Date int     `json:"date"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Type string  `json:"type"`
	}

	requestParams struct {
		Lang   string `query:"lang"`
		Name   string `query:"site_name"`
		Epoch  int    `query:"epoch_id" validate:"min=0,max=8"`
		Type   int    `query:"type_id" validate:"min=0,max=12"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}

	response map[string][]site
)

const (
	statement = `
    MATCH (s:Monument)<--(k:Knowledge)
    MATCH (s)-[:has]->(st:MonumentType)-[:translation {lang: {language}}]->(trType:Translate)
    MATCH (s)-[:has]->(e:Epoch)-[:translation {lang: {language}}]->(trEpoch:Translate)
    MATCH (r:Research)-[:has]->(k)
		WHERE %s
		WITH
			s.id as id,
			{
				site_name: collect(k.monument_name),
				research_name: collect(r.name),
				epoch: trEpoch.name,
				type: trType.name
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`

	entity = "Monument"
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	result, err := querySites(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response{
		"sites": result,
	})
}

func querySites(c echo.Context) ([]site, error) {
	neoHost := os.Getenv("Neo4jHost")
	DB, err := neoism.Connect(neoHost)
	assert.Nil(err)

	var res []site

	req := newRequestParams()
	if err = c.Bind(req); err != nil {
		return nil, routes.NotAllowedQueryParams
	}

	if err = c.Validate(req); err != nil {
		return nil, err
	}

	cq := neoism.CypherQuery{
		Statement: finalStatement(statement, siteFilterString(req)),
		Parameters: neoism.Props{
			"language": req.Lang,
			"name":     "(?ui).*" + req.Name + ".*$",
			"epoch":    req.Epoch,
			"type":     req.Type,
			"offset":   req.Offset,
			"limit":    req.Limit,
		},
		Result: &res,
	}

	err = DB.Cypher(&cq)
	assert.Nil(err)

	if len(res) > 0 {
		ids := make([]int, len(res))
		coords := make([]coordinates, len(res))

		for i, v := range res {
			ids[i] = v.ID
		}

		coordStatement := cypher.BuildCoordinates(ids, entity)
		coordsCQ := neoism.CypherQuery{
			Statement: coordStatement,
			Parameters: neoism.Props{
				"language": req.Lang,
			},
			Result: &coords,
		}

		err = DB.Cypher(&coordsCQ)
		assert.Nil(err)

		for i := range coords {
			res[i].Coords = coords[i]
		}

		return res, nil
	}

	return []site{}, nil
}

func siteFilterString(reqParams *requestParams) string {
	var filter []string

	if reqParams.Name != "" {
		filter = append(filter, "k.monument_name =~ {name}")
	}
	if reqParams.Epoch != 0 {
		filter = append(filter, "e.id = {epoch}")
	}
	if reqParams.Type != 0 {
		filter = append(filter, "st.id = {type}")
	}

	return strings.Join(filter, " AND ")
}

func finalStatement(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}

func newRequestParams() *requestParams {
	return &requestParams{
		Lang:   "en",
		Name:   "",
		Epoch:  0,
		Type:   0,
		Offset: 0,
		Limit:  20,
	}
}
