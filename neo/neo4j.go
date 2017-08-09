package neo

import (
	"os"

	"github.com/jmcvetta/neoism"
)

// DB is neo4j database instance
var DB *neoism.Database

// InitDB connecting to Neoj
func InitDB() (err error) {
	neoHost := os.Getenv("Neo4jHost")
	DB, err = neoism.Connect(neoHost)
	return err
}

// BuildCypherQuery return neoism library struct for quering Neo4j
func BuildCypherQuery(stmt string, dst interface{}, props neoism.Props) neoism.CypherQuery {
	return neoism.CypherQuery{
		Statement:  stmt,
		Result:     dst,
		Parameters: props,
	}
}
