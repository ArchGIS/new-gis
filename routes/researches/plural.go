package researches

import (
	"fmt"
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
		ResName    string `json:"research_name"`
		ReportName string `json:"report_name"`
		Year       int64  `json:"year"`
		Author     string `json:"author_name"`
		Type       string `json:"res_type"`
	}

	research struct {
		ID     uint64        `json:"id"`
		Item   item          `json:"item"`
		Coords []coordinates `json:"coordinates"`
	}

	resCoord struct {
		Rows []coordinates `json:"rows"`
	}
	coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	requestParams struct {
		Lang   string `query:"lang"`
		Name   string `query:"res_name"`
		Year   int64  `query:"res_year"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}
)

const (
	statement = `
		MATCH (n:Research)-[:hasauthor]->(a:Author)
		MATCH (n)-[:has]->(rep:Report)
		MATCH (n)-[:has]->(rtype:ResearchType)-[:translation {lang: {language}}]->(trans:Translate)
		WHERE %s
		WITH n.id as id, {
			research_name: n.name,
			report_name: rep.name,
			year: n.year,
			author_name: a.name,
			res_type: trans.name
		} as item
		RETURN id, item
		SKIP {offset} LIMIT {limit}
	`

	entity = "Research"
	minInt = -int64(1<<63 - 1)
)

// Plural gets info about researches
func Plural(c echo.Context) error {
	researches, err := queryResearches(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"researches": researches})
}

func queryResearches(c echo.Context) (researches []research, err error) {
	req := &requestParams{
		Lang:   "en",
		Name:   "",
		Year:   minInt,
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
		finalStatement(statement, researchFilterString(req)),
		&researches,
		neoism.Props{
			"language": req.Lang,
			"name":     neo.BuildRegexpFilter(req.Name),
			"year":     req.Year,
			"offset":   req.Offset,
			"limit":    req.Limit,
		},
	)

	err = neo.DB.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	if len(researches) > 0 {
		ids := make([]uint64, len(researches))
		for i, v := range researches {
			ids[i] = v.ID
		}

		listResearchSites := []cypher.ListIdsOfSites{}
		cq = neo.BuildCypherQuery(
			cypher.IdsOfResearchSites(ids),
			&listResearchSites,
			neoism.Props{},
		)
		err = neo.DB.Cypher(&cq)
		if err != nil {
			return nil, err
		}

		coords := make([]resCoord, len(researches))

		coordsCQ := neo.BuildCypherQuery(
			cypher.ResearchCoords(listResearchSites),
			&coords,
			neoism.Props{},
		)

		err = neo.DB.Cypher(&coordsCQ)
		if err != nil {
			return nil, err
		}

		for i := range coords {
			researches[i].Coords = coords[i].Rows
		}
	}

	return researches, nil
}

func researchFilterString(reqParams *requestParams) string {
	var filter []string

	if reqParams.Name != "" {
		filter = append(filter, "n.name =~ {name}")
	}
	if reqParams.Year != minInt {
		filter = append(filter, "n.year = {year}")
	}

	return strings.Join(filter, " AND ")
}

func finalStatement(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}
