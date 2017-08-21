package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

// Count returns count of entities in DB
func Count(c echo.Context) error {
	counts, err := Model.db.Counts()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"counts": counts})
}
