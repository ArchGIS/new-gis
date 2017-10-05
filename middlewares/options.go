package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AddOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Next()
	}
}

func HandleOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Allow", "OPTIONS, GET, POST")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			return
		}

		c.Next()
	}
}
