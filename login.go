package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	authSecret = "AUTH_SECRET"
)

type User struct {
	Name     string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

func loginHandler(c *gin.Context) {
	user := new(User)
	if err := c.Bind(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	if isAuthentificated(user.Name, user.Password) {
		// Fix this
		token := "token"
		duration := time.Now().Add(time.Hour * 24 * 90).Unix()

		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"expired": fmt.Sprintf("%d", duration),
		})
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func isAuthentificated(login, password string) bool {
	if login == "admin" && password == "qwerty" {
		return true
	}

	return false
}
