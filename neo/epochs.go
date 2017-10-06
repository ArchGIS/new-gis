package neo

import (
	"github.com/gin-gonic/gin"
)

type epoch struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Epochs(req gin.H) (epochs []epoch, err error) {
	// cq := BuildCypherQuery(epochsStatement, &epochs, neoism.Props{"language": req["lang"]})

	// err = db.Cypher(&cq)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

const (
	epochsStatement = `
		MATCH (n:Epoch)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)
