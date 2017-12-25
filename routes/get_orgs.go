package routes

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type (
	requestOrgs struct {
		Name string `form:"name"`
	}
)

func Organizations(c *gin.Context) {
	var req requestOrgs

	if err := c.Bind(&req); err != nil {
		log.Printf("bind didn't work: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		panic(err)
	}

	orgs, err := db.Organizations(map[string]interface{}{"name": req.Name})
	if err != nil {
		log.Printf("query to db failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"orgs": orgs})
}
