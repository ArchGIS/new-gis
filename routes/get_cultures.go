package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	requestCulture struct {
		Lang string `query:"lang"`
		Name string `query:"name"`
	}
)

func Cultures(c *gin.Context) {
	req := &requestCulture{
		Lang: "en",
		Name: "",
	}

	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	}

	cultures, err := db.Cultures(gin.H{
		"lang": req.Lang,
		"name": req.Name,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"cultures": cultures})
}
