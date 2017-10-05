package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CityRequest struct {
	Lang string `query:"lang"`
	Name string `query:"name"`
}

func Cities(c *gin.Context) {
	req := CityRequest{
		Lang: "en",
		// Name: "",
	}

	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	}

	cities, err := Model.db.Cities(gin.H{
		"lang": req.Lang,
		"name": req.Name,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"cities": cities})
}
