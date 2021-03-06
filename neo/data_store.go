package neo

import (
	"fmt"
	"strings"

	"database/sql"

	_ "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

type (
	DataStore interface {
		Graphql([]byte) (interface{}, error)

		RawQuery(string) (interface{}, error)

		Authors(map[string]interface{}) ([]authorsProps, error)
		Counts() ([]*nodesCounter, error)
		Cities(map[string]interface{}) ([]*cityProps, error)
		Cultures(map[string]interface{}) ([]*cultureProps, error)
		Epochs(map[string]interface{}) ([]epochProps, error)
		Organizations(map[string]interface{}) ([]*organisationProps, error)
		SiteTypes(map[string]interface{}) ([]siteTypeProps, error)

		Sites(map[string]interface{}) ([]*site, error)
		// Researches(gin.H) ([]pluralResearch, error)
	}

	DB struct {
		*sql.DB
	}
)

// InitDB connecting to Neoj
func InitDB(source string) (*DB, error) {
	db, err := sql.Open("neo4j-bolt", source)
	if err != nil {
		return nil, fmt.Errorf("can not open db connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can not initialize db with ping: %v", err)
	}

	return &DB{db}, nil
}

// addRegexpFilter return neo4j regexp filter
// for case-insensitive text search
func addRegexpFilter(par map[string]interface{}, keys ...string) {
	for _, v := range keys {
		par[v] = fmt.Sprintf("(?ui).*%s.*$", par[v].(string))
	}
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
