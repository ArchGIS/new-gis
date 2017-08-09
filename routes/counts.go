package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/neo"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type count struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

const (
	statement = `
		OPTIONAL MATCH (a:Author)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Research)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Heritage)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Monument)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Artifact)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Radiocarbon)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Report)
		RETURN labels(a)[0] as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Monography)
		RETURN "Monography" as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:Article)
		RETURN "Article" as name, count(a) as count
		UNION
		OPTIONAL MATCH (a:ArchiveDoc)
		RETURN "ArchiveDoc" as name, count(a) as count
	`
)

// Count returns count of entities in DB
func Count(c echo.Context) error {
	counts, err := queryCounts(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"counts": counts})
}

func queryCounts(c echo.Context) (counts []count, err error) {
	cq := neo.BuildCypherQuery(statement, &counts, neoism.Props{})

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return counts, nil
}
