package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/neo"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

const (
	epochsStatement = `
		MATCH (n:Epoch)-[:translation {lang: {language}}]->(tr:Translate)
		RETURN n.id as id, tr.name as name
	`
)

type (
	epoch struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

// Epochs return list of epochs
func Epochs(c echo.Context) (err error) {
	req := &request{Lang: "en"}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	var epochs []epoch

	cq := neo.BuildCypherQuery(epochsStatement, &epochs, neoism.Props{"language": req.Lang})

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"epochs": epochs})
}
