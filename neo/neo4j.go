package neo

import (
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/neoism"
)

type (
	DataStore interface {
		Counts() ([]NodesCounter, error)
		Cities(gin.H) ([]City, error)
		Cultures(gin.H) ([]Culture, error)
		Epochs(gin.H) ([]epoch, error)
		Organizations(gin.H) ([]organisation, error)
		SiteTypes(gin.H) ([]siteType, error)

		GetSite(string, string) (interface{}, error)
		QuerySiteResearches(id, lang string) (interface{}, error)
		Sites(gin.H) ([]pluralSite, error)
		Researches(gin.H) ([]pluralResearch, error)
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
