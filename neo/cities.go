package neo

import (
	"github.com/ArchGIS/new-gis/cypher"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type City struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Cities(req echo.Map) (cities []City, err error) {
	cq := BuildCypherQuery(
		cypher.Filter(cityStatement, filterCity(req)),
		&cities,
		neoism.Props{
			"language": req["lang"],
			"name":     BuildRegexpFilter(req["name"]),
		},
	)

	err = db.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func filterCity(req echo.Map) (filter string) {
	if req["name"] != "" {
		filter = "WHERE n.name =~ {name}"
	}

	return filter
}

const cityStatement = `
	MATCH (n:City)
	%s
	RETURN n.id as id, n.name as name
`
