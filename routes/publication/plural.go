package publication

import (
	"net/http"
	"strings"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	item struct {
		Name       string `json:"name"`
		AuthorName string `json:"author_name"`
		PubYear    int64  `json:"pub_year"`
	}

	publication struct {
		ID   uint64 `json:"id"`
		Item item   `json:"item"`
	}

	requestParams struct {
		Name   string `query:"name"`
		Year   int64  `query:"year"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}
)

const (
	statement = `
    MATCH (n:Publication)-->(a:Author)
		%s
		WITH
			n.id as id,
			{
				name: n.name,
				author_name: a.name,
				pub_year: n.published_at
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	pubs, err := queryPubs(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"pubs": pubs})
}

func queryPubs(c echo.Context) (pubs []publication, err error) {
	req := &requestParams{
		Name:   "",
		Year:   routes.MinInt,
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
		cypher.Filter(statement, pubFilterString(req)),
		&pubs,
		neoism.Props{
			"name":   neo.BuildRegexpFilter(req.Name),
			"year":   req.Year,
			"offset": req.Offset,
			"limit":  req.Limit,
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return pubs, nil
}

func pubFilterString(reqParams *requestParams) string {
	var filter []string
	var stmt string

	if reqParams.Name != "" {
		filter = append(filter, "n.name =~ {name}")
	}
	if reqParams.Year != routes.MinInt {
		filter = append(filter, "n.year = {year}")
	}
	if len(filter) > 0 {
		stmt += "WHERE " + strings.Join(filter, " AND ")
	}

	return stmt
}
