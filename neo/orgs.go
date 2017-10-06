package neo

import (
	"github.com/gin-gonic/gin"
)

type organisation struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Organizations(req gin.H) (orgs []organisation, err error) {
	// cq := BuildCypherQuery(
	// 	cypher.Filter(orgStatement, filterOrgs(req)),
	// 	&orgs,
	// 	neoism.Props{"name": cypher.BuildRegexpFilter(req["name"])},
	// )

	// err = db.Cypher(&cq)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func filterOrgs(req gin.H) (filter string) {
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
