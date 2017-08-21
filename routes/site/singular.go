package site

import (
	"net/http"

	"github.com/labstack/echo"
)

// Singular get info about single archaelogical site
func Singular(c echo.Context) error {
	id := c.Param("id")

	site, err := getSite(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"site": site})
}

func getSite(id string) (interface{}, error) {

	return nil, nil
}
