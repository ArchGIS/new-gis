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
