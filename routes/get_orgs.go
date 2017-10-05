package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	requestOrgs struct {
		Name string `query:"name"`
	}
)

func Organizations(c *gin.Context) {
	// req := &requestOrgs{Name: ""}
	var req requestOrgs

	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	}

	orgs, err := db.Organizations(gin.H{"name": req.Name})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"orgs": orgs})
}
