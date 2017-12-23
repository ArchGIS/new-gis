package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	requestAuthor struct {
		Name   string `form:"name"`
		Offset int    `form:"offset"`
		Limit  int    `form:"limit"`
	}
)

// QueryAuthors gets info about authors
func QueryAuthors(c *gin.Context) {
	req := requestAuthor{Limit: 10}

	if err := c.Bind(&req); err != nil {
		log.Printf("could not bind request: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		return
	}

	authors, err := db.Authors(map[string]interface{}{
		"name":   req.Name,
		"offset": req.Offset,
		"limit":  req.Limit,
	})
	if err != nil {
		log.Printf("db query failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"authors": authors})
}
