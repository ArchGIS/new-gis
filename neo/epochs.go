package neo

import (
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type epoch struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Epochs(req echo.Map) (epochs []epoch, err error) {
	cq := BuildCypherQuery(epochsStatement, &epochs, neoism.Props{"language": req["lang"]})

	err = db.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return epochs, nil
}

const (
	epochsStatement = `
		MATCH (n:Epoch)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)
