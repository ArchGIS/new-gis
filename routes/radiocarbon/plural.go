package radiocarbon

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
		Index string `json:"index"`
		Date  int    `json:"date"`
	}

	radiocarbon struct {
		ID     uint64      `json:"id"`
		Item   item        `json:"item"`
		Coords coordinates `json:"coordinates"`
	}

	coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	requestParams struct {
		Name   string `query:"index"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}
)

const (
	statement = `
		MATCH (rc:Radiocarbon)
		%s
		WITH
			rc.id as id,
			{
				index: rc.name,
				date: rc.date
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`
)

// Plural gets info about archeological sites
func Plural(c echo.Context) error {
	radiocarbon, err := queryRadiocarbon(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"radiocarbon": radiocarbon})
}

func queryRadiocarbon(c echo.Context) (rcs []radiocarbon, err error) {
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
		finalStatement(statement, carbonFilterString(req)),
		&rcs,
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

	rcLength := len(rcs)
	if rcLength > 0 {
		ids := make([]uint64, rcLength)
		coords := make([]coordinates, rcLength)

		for i, v := range rcs {
			ids[i] = v.ID
		}

		cq = neo.BuildCypherQuery(
			cypher.BuildCoordinates(ids, "Radiocarbon", false),
			&coords,
			neoism.Props{},
		)

		err = neo.DB.Cypher(&cq)
		if err != nil {
			return nil, err
		}

		for i, v := range coords {
			rcs[i].Coords = v
		}
	}

	return rcs, nil
}

func carbonFilterString(reqParams *requestParams) (filter string) {
	if reqParams.Name != "" {
		filter = "WHERE rc.name =~ {name}"
	}

	return filter
}

func finalStatement(statement, filter string) string {
	return fmt.Sprintf(statement, filter)
}
