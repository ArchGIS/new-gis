package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CityRequest struct {
	Lang string `form:"lang"`
	Name string `form:"name"`
}

func Cities(c *gin.Context) {
	req := CityRequest{Lang: "en"}

	if err := c.BindQuery(&req); err != nil {
		log.Printf("bind didn't work: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		panic(err)
	}

	cities, err := db.Cities(map[string]interface{}{
		"lang": req.Lang,
		"name": req.Name,
	})
	if err != nil {
		log.Printf("query to db failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"cities": cities})
}
