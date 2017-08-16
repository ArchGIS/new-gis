package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	requestCity struct {
		Lang string `query:"lang"`
		Name string `query:"name"`
	}

	city struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

const cityStatement = `
	MATCH (n:City)
	%s
	RETURN n.id as id, n.name as name
`

func Cities(c echo.Context) (err error) {
	req := &requestCity{
		Lang: "en",
		Name: "",
	}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	var cities []city
	cq := neo.BuildCypherQuery(
		cypher.Filter(cityStatement, filterCity(req)),
		&cities,
		neoism.Props{
			"language": req.Lang,
			"name":     neo.BuildRegexpFilter(req.Name),
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"cities": cities})
}

func filterCity(req *requestCity) (filter string) {
	if req.Name != "" {
		filter = "WHERE n.name =~ {name}"
	}

	return filter
}
