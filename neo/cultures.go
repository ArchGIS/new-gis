package neo

import (
	"github.com/gin-gonic/gin"
)

type cultureProps struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// func (db *DB) Cultures(req gin.H) ([]cultureProps, error) {
// 	rows, err := db.QueryNeo(
// 		fmt.Sprintf(cultureStatement, filterCulture(req)),
// 		gin.H{
// 			"language": req["lang"],
// 			"name":     buildRegexpFilter(req["name"]),
// 		},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	data, _, err := rows.All()
// 	if err != nil {
// 		return nil, err
// 	}

// 	cultures := make([]cultureProps, len(data))
// 	for i, row := range data {
// 		cultures[i] = cultureProps{
// 			ID:   row[0].(int64),
// 			Name: row[1].(string),
// 		}
// 	}

// 	return cultures, nil
// }

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
