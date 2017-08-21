package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	requestCulture struct {
		Lang string `query:"lang"`
		Name string `query:"name"`
	}
)

func Cultures(c echo.Context) (err error) {
	req := &requestCulture{
		Lang: "en",
		Name: "",
	}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	cultures, err := Model.db.Cultures(echo.Map{
		"lang": req.Lang,
		"name": req.Name,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"cultures": cultures})
}
