package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SiteTypes return list with types of sites
func SiteTypes(c *gin.Context) {
	req := request{Lang: "en"}

	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	}

	siteTypes, err := db.SiteTypes(gin.H{"lang": req.Lang})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"siteTypes": siteTypes})
}
