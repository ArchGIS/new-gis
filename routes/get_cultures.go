package routes

import (
	"fmt"
	"net/http"

	"github.com/ArchGIS/new-gis/neo"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	requestCulture struct {
		Lang string `query:"lang"`
		Name string `query:"name"`
	}

	culture struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

const cultureStatement = `
	MATCH (n:Culture)-[:translation {lang: {language}}]->(tr:Translate)
	%s
	RETURN n.id as id, tr.name as name
`

func Cultures(c echo.Context) (err error) {
	req := &requestCulture{
		Lang: "en",
		Name: "",
	}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	var cultures []culture
	cq := neo.BuildCypherQuery(
		finalStatement(cultureStatement, filterCulture(req)),
		&cultures,
		neoism.Props{
			"language": req.Lang,
			"name":     neo.BuildRegexpFilter(req.Name),
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"cultures": cultures})
}

func filterCulture(req *requestCulture) (filter string) {
	if req.Name != "" {
		filter = "WHERE tr.name =~ {name}"
	}

	return filter
}

func finalStatement(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}
