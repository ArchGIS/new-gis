package neo

import (
	"github.com/ArchGIS/new-gis/cypher"
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/neoism"
)

type Culture struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Cultures(req gin.H) (cultures []Culture, err error) {
	cq := BuildCypherQuery(
		cypher.Filter(cultureStatement, filterCulture(req)),
		&cultures,
		neoism.Props{
			"language": req["lang"],
			"name":     cypher.BuildRegexpFilter(req["name"]),
		},
	)

	err = db.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return cultures, nil
}

func filterCulture(req gin.H) (filter string) {
	if req["name"] != "" {
		filter = "WHERE tr.name =~ {name}"
	}

	return filter
}

const cultureStatement = `
	MATCH (n:Culture)-[:translation {lang: {language}}]->(tr:Translate)
	%s
	RETURN n.id as id, tr.name as name
`
