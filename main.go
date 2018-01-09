package main

import (
	"os"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	middle "github.com/ArchGIS/new-gis/middlewares"
	"github.com/ArchGIS/new-gis/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	err := routes.InitEnv(os.Getenv("NEO4J_BOLT"))
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(middle.HandleOptions())

	r.Use(static.Serve("/", static.LocalFile(os.Getenv("STATIC_PATH"), true)))

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(os.Getenv("AUTH_SECRET")),
		Timeout:    time.Hour * 24 * 90,
		MaxRefresh: time.Hour * 24 * 90,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if userId == "admin" && password == "qwerty" {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	}

	r.OPTIONS("/login", func(*gin.Context) {})
	r.POST("/login", authMiddleware.LoginHandler)

	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/refresh_token", authMiddleware.RefreshHandler)

		api.POST("/graphql", routes.Graphql)

		api.GET("/counts", routes.Count)
		api.GET("/epochs", routes.Epochs)
		api.GET("/site_types", routes.SiteTypes)
		api.GET("/cultures", routes.Cultures)
		api.GET("/cities", routes.Cities)
		api.GET("/organizations", routes.Organizations)

		api.GET("/sites", routes.Sites)
		api.GET("/authors", routes.QueryAuthors)

		api.POST("/example", routes.Example)
	}

	r.Run(":7000")
}
