package heritage

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
		Name    string `json:"name"`
		Address string `json:"address"`
		Date    string `json:"date"`
	}

	heritage struct {
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
    MATCH (h:Heritage)
		%s
		WITH
			h.id as id,
			{
				name: h.name,
				address: h.address,
				date: h.docDate
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	heritages, err := queryHeritages(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"heritages": heritages})
}

func queryHeritages(c echo.Context) (heritages []heritage, err error) {
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
		cypher.Filter(statement, heritageFilterString(req)),
		&heritages,
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

	heritageLength := len(heritages)
	if heritageLength > 0 {
		ids := make([]uint64, heritageLength)
		coords := make([]coordinates, heritageLength)

		for i, v := range heritages {
			ids[i] = v.ID
		}

		cq = neo.BuildCypherQuery(
			cypher.BuildCoordinates(ids, "Heritage", false),
			&coords,
			neoism.Props{},
		)

		err = neo.DB.Cypher(&cq)
		if err != nil {
			return nil, err
		}

		for i, v := range coords {
			heritages[i].Coords = v
		}
	}

	return heritages, nil
}

func heritageFilterString(reqParams *requestParams) (filter string) {
	if reqParams.Name != "" {
		filter = "WHERE h.name =~ {name}"
	}

	return filter
}
