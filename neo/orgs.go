package neo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/encoding"
)

type organisationProps struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Organizations(req map[string]interface{}) ([]*organisationProps, error) {
	stmt := fmt.Sprintf(orgStatement, filterOrgs(req))

	addRegexpFilter(req, "name")
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not encode to gob: %v", err)
	}

	rows, err := db.Query(stmt, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*organisationProps
	for rows.Next() {
		org := new(organisationProps)
		err = rows.Scan(&org.ID, &org.Name)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		orgs = append(orgs, org)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
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
