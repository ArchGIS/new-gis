package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Count returns count of entities in DB
func Count(c *gin.Context) {
	counts, err := db.Counts()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}
