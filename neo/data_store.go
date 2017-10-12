package neo

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

type (
	DataStore interface {
		Counts() ([]nodesCounter, error)
		Cities(gin.H) ([]cityProps, error)
		Cultures(gin.H) ([]cultureProps, error)
		Epochs(gin.H) ([]epochProps, error)
		Organizations(gin.H) ([]organisationProps, error)
		SiteTypes(gin.H) ([]siteTypeProps, error)

		Sites(gin.H) ([]pluralSite, error)
		Researches(gin.H) ([]pluralResearch, error)

		GetSite(int64, string) (interface{}, error)
		QuerySiteResearches(id, lang string) (interface{}, error)
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

// buildRegexpFilter return neo4j regexp filter
// for case-insensitive text search
func buildRegexpFilter(needle interface{}) string {
	return fmt.Sprintf("(?ui).*%s.*$", needle)
}

// BuildCoordinates generates cypher query for searching actual coordinate
func actualCoordinates(length int, entity string) string {
	actualCoordQuery := `
		OPTIONAL MATCH (s:%s {id: {id%d}})-[:has]->(sr:SpatialReference)-[:has]->(srt:SpatialReferenceType)
		WITH sr, srt
		ORDER BY srt.id ASC, sr.date DESC LIMIT 1
		RETURN sr.x as x, sr.y as y
	`
	queries := make([]string, length)

	queries[0] = fmt.Sprintf(actualCoordQuery, entity, 0)

	for i := 1; i < length; i++ {
		queries[i] = fmt.Sprintf(actualCoordQuery, entity, i)
	}

	return strings.Join(queries, " UNION ALL ")
}
