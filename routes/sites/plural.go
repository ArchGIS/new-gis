package sites

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ArchGIS/new-gis/assert"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type monument struct {
	// `json:` tags matches column names in query
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ResName  string `json:"resName"`
	Epoch    string `json:"epoch"`
	SiteType string `json:"type"`
}

const (
	statement = `
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
		SKIP {offset} LIMIT {limit}
	`
)

func Plural(c echo.Context) error {
	result, err := querySites(c)

	if err == nil {
		return c.JSON(http.StatusOK, result)
	}

	return err
}

func querySites(c echo.Context) ([]monument, error) {
	neoHost := os.Getenv("Neo4jHost")
	DB, err := neoism.Connect(neoHost)
	assert.Nil(err)

	var res []monument

	name := c.QueryParam("site_name")

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	cq := neoism.CypherQuery{
		Statement: statement,
		Parameters: neoism.Props{
			"name":   "(?ui).*" + name + ".*$",
			"offset": offset,
			"limit":  limit,
		},
		Result: &res,
	}

	err = DB.Cypher(&cq)
	assert.Nil(err)

	if len(res) > 0 {
		return res, nil
	}

	return nil, echo.ErrNotFound
}
