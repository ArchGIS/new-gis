package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	authSecret = "AUTH_SECRET"
)

type User struct {
	Name     string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

func loginHandler(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return echo.ErrNotFound
	}

	if isAuthentificated(user.Name, user.Password) {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Admin"
		claims["admin"] = true
		duration := time.Now().Add(time.Hour * 24).Unix()
		claims["exp"] = duration

		// Generate encoded token
		t, err := token.SignedString([]byte(os.Getenv(authSecret)))
		if err != nil {
			return err
		}

		// Send it as response
		strExpired := fmt.Sprintf("%d", duration)
		return c.JSON(http.StatusOK, map[string]string{
			"token":   t,
			"expired": strExpired,
		})
	}

	return echo.ErrUnauthorized
}

func isAuthentificated(login, password string) bool {
	if login == "admin" && password == "qwerty" {
		return true
	}

	return false
}
