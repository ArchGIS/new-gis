package artifact

import (
	"fmt"
	"net/http"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	item struct {
		Name        string `json:"name"`
		Category    int    `json:"category"`
		Year        int64  `json:"year"`
		Collections string `json:"collections"`
	}

	artifact struct {
		ID     uint64      `json:"id"`
		Item   item        `json:"item"`
		Coords coordinates `json:"coordinates"`
	}

	coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	requestParams struct {
		Name   string `query:"name"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}
)

const (
	statement = `
		MATCH (a:Artifact)
		%s
		OPTIONAL MATCH (a)-->(ac:ArtifactCategory)
		OPTIONAL MATCH (a)-->(:StorageInterval)<--(coll:Collection)
		WITH
			a.id as id,
			{
				name: a.name,
				category: ac.id,
				collections: coll.name,
				year: a.year
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	artifacts, err := queryArtifacts(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"artifacts": artifacts})
}

func queryArtifacts(c echo.Context) (artifacts []artifact, err error) {
	req := &requestParams{
		Name:   "",
		Offset: 0,
		Limit:  20,
	}

	if err = c.Bind(req); err != nil {
		return nil, routes.NotAllowedQueryParams
	}

	if err = c.Validate(req); err != nil {
		return nil, routes.NotValidQueryParameters
	}

	cq := neo.BuildCypherQuery(
		finalStatement(statement, artifactFilterString(req)),
		&artifacts,
		neoism.Props{
			"name":   neo.BuildRegexpFilter(req.Name),
			"offset": req.Offset,
			"limit":  req.Limit,
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	artiLength := len(artifacts)
	if artiLength > 0 {
		ids := make([]uint64, artiLength)
		coords := make([]coordinates, artiLength)

		for i, v := range artifacts {
			ids[i] = v.ID
		}

		cq = neo.BuildCypherQuery(
			cypher.BuildCoordinates(ids, "Artifact", false),
			&coords,
			neoism.Props{},
		)

		err = neo.DB.Cypher(&cq)
		if err != nil {
			return nil, err
		}

		for i, v := range coords {
			artifacts[i].Coords = v
		}
	}

	return artifacts, nil
}

func artifactFilterString(reqParams *requestParams) (filter string) {
	if reqParams.Name != "" {
		filter = "WHERE a.name =~ {name}"
	}

	return filter
}

func finalStatement(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}
