package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/assert"
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

func Count(c echo.Context) error {
	result, err := queryCounts(c)

	if err == nil {
		return c.JSON(http.StatusOK, map[string][]count{
			"counts": result,
		})
	}

	return err
}

func queryCounts(c echo.Context) ([]count, error) {
	var err error
	var res []count

	cq := neoism.CypherQuery{
		Statement:  statement,
		Parameters: neoism.Props{},
		Result:     &res,
	}

	err = neo.DB.Cypher(&cq)
	assert.Nil(err)

	if len(res) > 0 {
		return res, nil
	}

	return nil, echo.ErrNotFound
}
