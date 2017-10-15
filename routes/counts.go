package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Count returns count of entities in DB
func Count(c *gin.Context) {
	counts, err := db.Counts()
	if err != nil {
		log.Printf("query to db failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}
