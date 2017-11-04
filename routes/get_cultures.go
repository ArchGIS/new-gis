package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	requestCulture struct {
		Lang string `form:"lang"`
		Name string `form:"name"`
	}
)

func Cultures(c *gin.Context) {
	req := requestCulture{Lang: "ru"}

	if err := c.Bind(&req); err != nil {
		log.Printf("could not bind request: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		return
	}

	cultures, err := db.Cultures(map[string]interface{}{
		"lang": req.Lang,
		"name": req.Name,
	})
	if err != nil {
		log.Printf("db query failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"cultures": cultures})
}
