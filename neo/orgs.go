package neo

import (
	"github.com/ArchGIS/new-gis/cypher"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type organisation struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Organizations(req echo.Map) (orgs []organisation, err error) {
	cq := BuildCypherQuery(
		cypher.Filter(orgStatement, filterOrgs(req)),
		&orgs,
		neoism.Props{"name": cypher.BuildRegexpFilter(req["name"])},
	)

	err = db.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func filterOrgs(req echo.Map) (filter string) {
	if req["name"] != "" {
		filter = "WHERE n.name =~ {name}"
	}

	return filter
}

const orgStatement = `
	MATCH (n:Organisation)
	%s
	RETURN n.id as id, n.name as name
`
