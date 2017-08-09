package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/neo"
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

// SiteTypes return list with types of sites
func SiteTypes(c echo.Context) (err error) {
	req := &request{Lang: "en"}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	var siteTypes []siteType

	cq := neo.BuildCypherQuery(siteTypesStatement, &siteTypes, neoism.Props{"language": req.Lang})

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"siteTypes": siteTypes})
}
