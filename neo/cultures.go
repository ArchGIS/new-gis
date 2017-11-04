package neo

import (
	"fmt"

	"github.com/johnnadratowski/golang-neo4j-bolt-driver/encoding"
)

type cultureProps struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Cultures(req map[string]interface{}) ([]*cultureProps, error) {
	stmt := fmt.Sprintf(cultureStatement, filterCulture(req))

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

	cultures := make([]*cultureProps, 0)
	for rows.Next() {
		cult := new(cultureProps)
		err = rows.Scan(&cult.ID, &cult.Name)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		cultures = append(cultures, cult)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return cultures, nil
}

func filterCulture(req map[string]interface{}) (filter string) {
	if req["name"] != "" {
		filter = "WHERE n[$lang + '_name'] =~ $name"
	}

	return filter
}

const cultureStatement = `
	MATCH (n:Culture)
	%s
	RETURN n.id as id, n[$lang + '_name'] as name
`
