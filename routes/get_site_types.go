package routes

import (
	"net/http"
	"os"

	"github.com/ArchGIS/new-gis/assert"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

const (
	siteTypesStatement = `
		MATCH (n:MonumentType)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)

type (
	siteType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

// Epochs return list of epochs
func SiteTypes(c echo.Context) error {
	neoHost := os.Getenv("Neo4jHost")
	DB, err := neoism.Connect(neoHost)
	assert.Nil(err)

	req := &request{Lang: "en"}
	if err = c.Bind(req); err != nil {
		return err
	}

	var res []siteType

	cq := neoism.CypherQuery{
		Statement: siteTypesStatement,
		Parameters: neoism.Props{
			"language": req.Lang,
		},
		Result: &res,
	}

	err = DB.Cypher(&cq)
	assert.Nil(err)

	if len(res) > 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"siteTypes": res,
		})
	}

	return echo.ErrNotFound
}
