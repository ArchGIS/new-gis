package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/neo"
	"github.com/labstack/echo"
)

type resRequestParams struct {
	Lang   string `query:"lang"`
	Name   string `query:"res_name"`
	Year   int64  `query:"res_year"`
	Offset int    `query:"offset"`
	Limit  int    `query:"limit"`
}

func Researches(c echo.Context) (err error) {
	req := &resRequestParams{
		Lang:   "en",
		Name:   "",
		Year:   neo.MinInt,
		Offset: 0,
		Limit:  20,
	}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	if err = c.Validate(req); err != nil {
		return NotValidQueryParameters
	}

	res, err := Model.db.Researches(echo.Map{
		"lang":   req.Lang,
		"name":   req.Name,
		"year":   req.Year,
		"offset": req.Offset,
		"limit":  req.Limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"researches": res})
}
