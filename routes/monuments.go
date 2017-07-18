package routes

import (
	"net/http"
	"os"

	"github.com/ArchGIS/ArchGoGIS/dbg"
	"github.com/ArchGIS/new-gis/assert"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

var DB *neoism.Database

type monument struct {
	// `json:` tags matches column names in query
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ResName  string `json:"resName"`
	Epoch    string `json:"epoch"`
	SiteType string `json:"type"`
}

func Monuments(c echo.Context) error {
	var err error
	neoHost := os.Getenv("Neo4jHost")
	DB, err = neoism.Connect(neoHost)
	assert.Nil(err)
	result, err := queryMonuments(c)

	if err == nil {
		return c.JSON(http.StatusOK, result)
	}

	return err
}

func queryMonuments(c echo.Context) ([]monument, error) {
	var res []monument

	cq := neoism.CypherQuery{
		Statement: `
			MATCH (s:Monument)<--(k:Knowledge)
			WHERE k.monument_name =~ {name}
			with s, k
			MATCH (s)-[:has]->(st:MonumentType)
			MATCH (s)-[:has]->(e:Epoch)
			MATCH (r:Research)-[:has]->(k)
			RETURN
				s.id as id,
				k.monument_name as name,
				r.name as resName,
				e.name as epoch,
				st.name as type
		`,
		Parameters: neoism.Props{"name": "(?ui).*болг.*$"},
		Result:     &res,
	}

	err := DB.Cypher(&cq)
	assert.Nil(err)

	dbg.Dump(res)
	if len(res) > 0 {
		return res, nil
	}

	return nil, echo.ErrNotFound
}
