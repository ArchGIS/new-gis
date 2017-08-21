package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

type CityRequest struct {
	Lang string `query:"lang"`
	Name string `query:"name"`
}

func Cities(c echo.Context) (err error) {
	req := &CityRequest{
		Lang: "en",
		Name: "",
	}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	cities, err := Model.db.Cities(echo.Map{
		"lang": req.Lang,
		"name": req.Name,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"cities": cities})
}
