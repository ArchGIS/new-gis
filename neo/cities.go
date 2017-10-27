package neo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/encoding"
)

type cityProps struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Cities(req map[string]interface{}) ([]*cityProps, error) {
	stmt := fmt.Sprintf(cityStatement, filterCity(req))

	// req["name"] = `(?ui).*` + req["name"].(string) + `+.*`
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

	cities := make([]*cityProps, 0)
	for rows.Next() {
		city := new(cityProps)
		err = rows.Scan(&city.ID, &city.Name)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		cities = append(cities, city)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return cities, nil
}

func filterCity(req gin.H) (filter string) {
	if req["name"] != "" {
		filter = `WHERE n.name =~ {name}`
	}

	return filter
}

const cityStatement = `
	MATCH (n:City)
	%s
	RETURN n.id as id, n.name as name
`
