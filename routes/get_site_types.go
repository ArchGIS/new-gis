package routes

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

// SiteTypes return list with types of sites
func SiteTypes(c *gin.Context) {
	req := request{Lang: "en"}

	if err := c.Bind(&req); err != nil {
		log.Printf("could not bind request: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		return
	}

	siteTypes, err := db.SiteTypes(map[string]interface{}{"lang": req.Lang})
	if err != nil {
		log.Printf("db query failed:: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		return
	}

	c.JSON(http.StatusOK, gin.H{"siteTypes": siteTypes})
}
