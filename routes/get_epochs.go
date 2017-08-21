package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

// Epochs return list of epochs
func Epochs(c echo.Context) (err error) {
	req := &request{Lang: "en"}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	epochs, err := Model.db.Epochs(echo.Map{"lang": req.Lang})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"epochs": epochs})
}
