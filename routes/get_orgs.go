package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	requestOrgs struct {
		Name string `query:"name"`
	}

	organisation struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

const orgStatement = `
	MATCH (n:Organisation)
	%s
	RETURN n.id as id, n.name as name
`

func Organizations(c echo.Context) (err error) {
	req := &requestOrgs{Name: ""}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	var orgs []organisation
	cq := neo.BuildCypherQuery(
		cypher.Filter(orgStatement, filterOrgs(req)),
		&orgs,
		neoism.Props{"name": neo.BuildRegexpFilter(req.Name)},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"orgs": orgs})
}

func filterOrgs(req *requestOrgs) (filter string) {
	if req.Name != "" {
		filter = "WHERE n.name =~ {name}"
	}

	return filter
}
