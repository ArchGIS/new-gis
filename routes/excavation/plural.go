package excavation

import (
	"net/http"
	"strings"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/neoism"
)

type (
	item struct {
		Name     string   `json:"exc_name"`
		ResName  []string `json:"res_name"`
		SiteName []string `json:"site_name"`
		Boss     string   `json:"boss"`
		Area     int      `json:"area"`
	}

	excavation struct {
		ID     uint64      `json:"id"`
		Item   item        `json:"item"`
		Coords coordinates `json:"coordinates"`
	}

	coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	requestParams struct {
		Name    string `query:"exc_name"`
		Author  string `query:"author_name"`
		ResYear int64  `query:"res_year"`
		Offset  int    `query:"offset"`
		Limit   int    `query:"limit"`
	}
)

const (
	statement = `
		MATCH (a:Author)<-[:hasauthor]-(r:Research)-->(k:Knowledge)-->(m:Monument)-->(e:Excavation)
		MATCH (r)-->(e)
		%s
		WITH
			e.id as id,
			{
				exc_name: e.name,
				res_name: COLLECT(r.name),
				site_name: COLLECT(k.monument_name),
				boss: e.boss,
				area: e.area
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`
)

// Plural gets info about archeological sites
func Plural(c *gin.Context) error {
	excavations, err := queryExcavations(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"excavations": excavations})
}

func queryExcavations(c *gin.Context) (excs []excavation, err error) {
	req := &requestParams{
		Name:    "",
		Author:  "",
		ResYear: routes.MinInt,
		Offset:  0,
		Limit:   20,
	}

	if err = c.Bind(req); err != nil {
		return nil, routes.NotAllowedQueryParams
	}

	if err = c.Validate(req); err != nil {
		return nil, routes.NotValidQueryParameters
	}

	cq := neo.BuildCypherQuery(
		cypher.Filter(statement, excFilterString(req)),
		&excs,
		neoism.Props{
			"name":     cypher.BuildRegexpFilter(req.Name),
			"author":   cypher.BuildRegexpFilter(req.Author),
			"res_year": req.ResYear,
			"offset":   req.Offset,
			"limit":    req.Limit,
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	excsLength := len(excs)
	if excsLength > 0 {
		ids := make([]uint64, excsLength)
		coords := make([]coordinates, excsLength)

		for i, v := range excs {
			ids[i] = v.ID
		}

		cq = neo.BuildCypherQuery(
			cypher.BuildCoordinates(ids, "Excavation", false),
			&coords,
			neoism.Props{},
		)

		err = neo.DB.Cypher(&cq)
		if err != nil {
			return nil, err
		}

		for i, v := range coords {
			excs[i].Coords = v
		}
	}

	return excs, nil
}

func excFilterString(reqParams *requestParams) string {
	var filter []string
	var stmt string

	if reqParams.Name != "" {
		filter = append(filter, "e.name =~ {name}")
	}
	if reqParams.Author != "" {
		filter = append(filter, "a.name =~ {author}")
	}
	if reqParams.ResYear != routes.MinInt {
		filter = append(filter, "r.year = {res_year}")
	}
	if len(filter) > 0 {
		stmt = "WHERE " + strings.Join(filter, " AND ")
	}

	return stmt
}
