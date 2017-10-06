package neo

import (
	"github.com/gin-gonic/gin"
)

type epoch struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Epochs(req gin.H) ([]epoch, error) {
	rows, err := db.QueryNeo(
		epochsStatement,
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

	epochs := make([]epoch, len(data))
	for i, row := range data {
		epochs[i] = epoch{
			ID:   row[0].(int64),
			Name: row[1].(string),
		}
	}

	return epochs, nil
}

const (
	epochsStatement = `
		MATCH (n:Epoch)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)
