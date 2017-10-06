package neo

import (
	"github.com/gin-gonic/gin"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
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
		bolt.Conn
	}
)

// InitDB connecting to Neoj
func InitDB(source string) (*DB, error) {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(source)
	if err != nil {
		return nil, err
	}
	return &DB{conn}, nil
}
