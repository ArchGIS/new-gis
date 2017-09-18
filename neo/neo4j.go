package neo

import (
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	DataStore interface {
		Counts() ([]NodesCounter, error)
		Cities(echo.Map) ([]City, error)
		Cultures(echo.Map) ([]Culture, error)
		Epochs(echo.Map) ([]epoch, error)
		Organizations(echo.Map) ([]organisation, error)
		SiteTypes(echo.Map) ([]siteType, error)

		GetSite(string, string) ([]knowledge, error)
		Sites(echo.Map) ([]pluralSite, error)
		Researches(echo.Map) ([]pluralResearch, error)
	}

	DB struct {
		*neoism.Database
	}
)

// InitDB connecting to Neoj
func InitDB(source string) (*DB, error) {
	// neoHost := os.Getenv("Neo4jHost")
	db, err := neoism.Connect(source)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// BuildCypherQuery return neoism library struct for quering Neo4j
func BuildCypherQuery(stmt string, dst interface{}, props neoism.Props) neoism.CypherQuery {
	return neoism.CypherQuery{
		Statement:  stmt,
		Result:     dst,
		Parameters: props,
	}
}
