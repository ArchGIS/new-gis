package neo

import (
	"github.com/gin-gonic/gin"
)

type (
	siteTypeProps struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
)

func (db *DB) SiteTypes(req gin.H) ([]siteTypeProps, error) {
	rows, err := db.QueryNeo(
		siteTypesStatement,
		gin.H{"language": req["lang"]},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, _, err := rows.All()
	if err != nil {
		return nil, err
	}

	sTypes := make([]siteTypeProps, len(data))
	for i, row := range data {
		sTypes[i] = siteTypeProps{
			ID:   row[0].(int64),
			Name: row[1].(string),
		}
	}

	return sTypes, nil
}

const (
	siteTypesStatement = `
		MATCH (n:MonumentType)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)
