package neo

import (
	"os"

	"github.com/jmcvetta/neoism"
)

var DB *neoism.Database

func InitDB() error {
	neoHost := os.Getenv("Neo4jHost")
	var err error
	DB, err = neoism.Connect(neoHost)
	return err
}
