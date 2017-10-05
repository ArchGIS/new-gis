package neo

import (
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/neoism"
)

type (
	siteType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

func (db *DB) SiteTypes(req gin.H) (siteTypes []siteType, err error) {
	cq := BuildCypherQuery(siteTypesStatement, &siteTypes, neoism.Props{"language": req["lang"]})

	err = db.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return siteTypes, nil
}

const (
	siteTypesStatement = `
		MATCH (n:MonumentType)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)
