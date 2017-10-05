package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Epochs return list of epochs
func Epochs(c *gin.Context) {
	req := request{Lang: "en"}

	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	}

	epochs, err := db.Epochs(gin.H{"lang": req.Lang})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"epochs": epochs})
}
