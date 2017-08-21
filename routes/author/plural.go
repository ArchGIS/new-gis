package author

import (
	"net/http"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	item struct {
		Name    string   `json:"author_name"`
		ResName []string `json:"research_name"`
	}

	author struct {
		ID   uint64 `json:"id"`
		Item item   `json:"item"`
	}

	requestParams struct {
		Name   string `query:"name"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}
)

const (
	statement = `
    MATCH (a:Author)<-[:hasauthor]-(r:Research)
		%s
		WITH
			a.id as id,
			{
				author_name: a.name,
				research_name: collect(r.name)
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	authors, err := queryAuthors(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"authors": authors})
}

func queryAuthors(c echo.Context) (authors []author, err error) {
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
		cypher.Filter(statement, authorFilterString(req)),
		&authors,
		neoism.Props{
			"name":   cypher.BuildRegexpFilter(req.Name),
			"offset": req.Offset,
			"limit":  req.Limit,
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	return authors, nil
}

func authorFilterString(reqParams *requestParams) (filter string) {
	if reqParams.Name != "" {
		filter = "WHERE a.name =~ {name}"
	}

	return filter
}
