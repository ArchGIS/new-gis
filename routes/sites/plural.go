package sites

import (
	"net/http"
	"os"

	"github.com/ArchGIS/new-gis/assert"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	item struct {
		Name     string `json:"site_name"`
		ResName  string `json:"research_name"`
		Epoch    string `json:"epoch"`
		SiteType string `json:"type"`
	}

	site struct {
		ID   int  `json:"id"`
		Item item `json:"item"`
	}

	requestParams struct {
		Name   string `query:"site_name"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}

	response map[string][]site
)

const (
	statement = `
		MATCH (s:Monument)<--(k:Knowledge)
		WHERE k.monument_name =~ {name}
		WITH s, k
		MATCH (s)-[:has]->(st:MonumentType)
		MATCH (s)-[:has]->(e:Epoch)
		MATCH (r:Research)-[:has]->(k)
		RETURN
			s.id AS id,
			{
				site_name: k.monument_name,
				research_name: r.name,
				epoch: e.name,
				type: st.name
			} as item
		SKIP {offset} LIMIT {limit}
	`
)

func Plural(c echo.Context) error {
	result, err := querySites(c)

	if err == nil {
		return c.JSON(http.StatusOK, response{
			"sites": result,
		})
	}

	return err
}

func querySites(c echo.Context) ([]site, error) {
	neoHost := os.Getenv("Neo4jHost")
	DB, err := neoism.Connect(neoHost)
	assert.Nil(err)

	var res []site

	req := new(requestParams)
	if err = c.Bind(req); err != nil {
		return nil, routes.NotAllowedQueryParams
	}

	cq := neoism.CypherQuery{
		Statement: statement,
		Parameters: neoism.Props{
			"name":   "(?ui).*" + req.Name + ".*$",
			"offset": req.Offset,
			"limit":  req.Limit,
		},
		Result: &res,
	}

	err = DB.Cypher(&cq)
	assert.Nil(err)

	if len(res) > 0 {
		return res, nil
	}

	return []site{}, nil
}
