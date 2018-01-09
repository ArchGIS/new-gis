package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Example(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}