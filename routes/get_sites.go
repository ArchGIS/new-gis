package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	requestParams struct {
		Name   string `query:"site_name"`
		Epoch  int    `query:"epoch_id" validate:"min=0,max=8"`
		Type   int    `query:"type_id" validate:"min=0,max=12"`
		Offset int    `query:"offset"`
		Limit  int    `query:"limit"`
	}
)

// Sites gets info about archeological sites
func Sites(c echo.Context) (err error) {
	req := &requestParams{
		Name:   "",
		Epoch:  0,
		Type:   0,
		Offset: 0,
		Limit:  20,
	}

	if err = c.Bind(req); err != nil {
		return NotAllowedQueryParams
	}

	if err = c.Validate(req); err != nil {
		return NotValidQueryParameters
	}

	sites, err := Model.db.Sites(echo.Map{
		"name":   req.Name,
		"epoch":  req.Epoch,
		"type":   req.Type,
		"offset": req.Offset,
		"limit":  req.Limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"sites": sites})
}

// SingleSite get info about single archaelogical site
func SingleSite(c echo.Context) error {
	id := c.Param("id")
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	site, err := Model.db.GetSite(id, lang)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"site": site})
}
