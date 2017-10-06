package neo

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type organisationProps struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Organizations(req gin.H) ([]organisationProps, error) {
	rows, err := db.QueryNeo(
		fmt.Sprintf(orgStatement, filterOrgs(req)),
		gin.H{"name": req["name"]},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, _, err := rows.All()
	if err != nil {
		return nil, err
	}

	orgs := make([]organisationProps, len(data))
	for i, row := range data {
		orgs[i] = organisationProps{
			ID:   row[0].(int64),
			Name: row[1].(string),
		}
	}

	return orgs, nil
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
