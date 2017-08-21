package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

// SiteTypes return list with types of sites
func SiteTypes(c echo.Context) (err error) {
	req := &request{Lang: "en"}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	siteTypes, err := Model.db.SiteTypes(echo.Map{"lang": req.Lang})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"siteTypes": siteTypes})
}
