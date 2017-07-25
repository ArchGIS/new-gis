package sites

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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
		Epoch  int    `query:"epoch_id"`
		Type   int    `query:"type_id"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}

	response map[string][]site
)

const (
	statement = `
		MATCH (s:Monument)<--(k:Knowledge)
		MATCH (s)-[:has]->(st:MonumentType)
		MATCH (s)-[:has]->(e:Epoch)
		MATCH (r:Research)-[:has]->(k)
		WHERE %s
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

	req := newRequestParams()
	if err = c.Bind(req); err != nil {
		return nil, routes.NotAllowedQueryParams
	}

	cq := neoism.CypherQuery{
		Statement: finalStatement(statement, siteFilterString(req)),
		Parameters: neoism.Props{
			"name":   "(?ui).*" + req.Name + ".*$",
			"epoch":  req.Epoch,
			"type":   req.Type,
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
		Name:   "",
		Epoch:  0,
		Type:   0,
		Offset: 0,
		Limit:  20,
	}
}
