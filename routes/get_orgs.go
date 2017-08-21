package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	requestOrgs struct {
		Name string `query:"name"`
	}
)

func Organizations(c echo.Context) (err error) {
	req := &requestOrgs{Name: ""}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	orgs, err := Model.db.Organizations(echo.Map{"name": req.Name})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"orgs": orgs})
}
