package routes

import (
	"net/http"

	"github.com/ArchGIS/new-gis/neo"
	"github.com/gin-gonic/gin"
)

type resRequestParams struct {
	Lang   string `query:"lang"`
	Name   string `query:"res_name"`
	Year   int64  `query:"res_year"`
	Offset int    `query:"offset"`
	Limit  int    `query:"limit"`
}

func Researches(c *gin.Context) {
	req := resRequestParams{
		Lang: "en",
		// Name:   "",
		Year: neo.MinInt,
		// Offset: 0,
		Limit: 20,
	}

	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	}

	// if err = c.Validate(req); err != nil {
	// 	return NotValidQueryParameters
	// }

	res, err := Model.db.Researches(gin.H{
		"lang":   req.Lang,
		"name":   req.Name,
		"year":   req.Year,
		"offset": req.Offset,
		"limit":  req.Limit,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"researches": res})
}
