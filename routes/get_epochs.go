package routes

import (
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
)

// Epochs return list of epochs
func Epochs(c *gin.Context) {
	req := request{Lang: "en"}

	if err := c.Bind(&req); err != nil {
		log.Printf("could not bind request: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		return
	}

	epochs, err := db.Epochs(map[string]interface{}{"lang": req.Lang})
	if err != nil {
		log.Printf("db query failed:: %v", err)
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
		return
	}

	c.JSON(http.StatusOK, gin.H{"epochs": epochs})
}
