package report

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ArchGIS/new-gis/dbg"
	"github.com/ArchGIS/new-gis/neo"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

type (
	item struct {
		ReportName   string `json:"report_name"`
		ResearchName string `json:"research_name"`
		AuthorName   string `json:"author_name"`
		ReportYear   int64  `json:"report_year"`
	}

	report struct {
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
    MATCH (a:Author)<-[:hasauthor]-(rep:Report)<-[:has]-(res:Research)
		%s
		WITH
			rep.id as id,
			{
				report_name: rep.name,
				research_name: res.name,
				author_name: a.name,
				report_year: rep.year
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	reports, err := queryReports(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"reports": reports})
}

func queryReports(c echo.Context) (reports []report, err error) {
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

	tmp := finalStatement(statement, reportFilterString(req))
	dbg.Dump(tmp)
	cq := neo.BuildCypherQuery(
		tmp,
		&reports,
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

	return reports, nil
}

func reportFilterString(reqParams *requestParams) string {
	var filter []string
	var stmt string

	if reqParams.Name != "" {
		filter = append(filter, "rep.name =~ {name}")
	}
	if reqParams.Year != routes.MinInt {
		filter = append(filter, "rep.year = {year}")
	}
	if len(filter) > 0 {
		stmt += "WHERE " + strings.Join(filter, " AND ")
	}

	return stmt
}

func finalStatement(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}
