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

	r.Use(middle.AddOrigin())
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

	r.Any("/login", authMiddleware.LoginHandler)

	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/refresh_token", authMiddleware.RefreshHandler)

		api.POST("/graphql", routes.Graphql)

		// apiV1.GET("/counts", routes.Count)
		api.GET("/epochs", routes.Epochs)
		api.GET("/site_types", routes.SiteTypes)
		api.GET("/cultures", routes.Cultures)
		api.GET("/cities", routes.Cities)
		api.GET("/organizations", routes.Organizations)

		api.GET("/sites", routes.Sites)
		// apiV1.GET("/researches", research.Plural)
		// api.GET("/authors", routes.QueryAuthors)
		// apiV1.GET("/reports", report.Plural)
		// apiV1.GET("/heritages", heritage.Plural)
		// apiV1.GET("/excavations", excavation.Plural)
		// apiV1.GET("/radiocarbons", radiocarbon.Plural)
		// apiV1.GET("/artifacts", artifact.Plural)
		// apiV1.GET("/publications", publication.Plural)
	}

	r.Run(":8181")
}
