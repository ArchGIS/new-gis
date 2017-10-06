package neo

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type City struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Cities(req gin.H) ([]City, error) {
	rows, err := db.QueryNeo(
		fmt.Sprintf(cityStatement, filterCity(req)),
		gin.H{
			"language": req["lang"],
			"name":     buildRegexpFilter(req["name"]),
		},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, _, err := rows.All()
	if err != nil {
		return nil, err
	}

	cities := make([]City, len(data))
	for i, row := range data {
		cities[i] = City{
			ID:   row[0].(int64),
			Name: row[1].(string),
		}
	}

	return cities, nil
}

func filterCity(req gin.H) (filter string) {
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
