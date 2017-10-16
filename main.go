package main

import (
	"os"

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
	err := routes.InitEnv(os.Getenv("Neo4jBolt"))
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// e.Validator = &CustomValidator{validator: validator.New()}

	r.Use(middle.AddOrigin())
	r.Use(middle.HandleOptions())

	r.Any("/login", loginHandler)

	apiV1 := r.Group("/v1")
	// apiV1.Use(middleware.JWT([]byte(os.Getenv(authSecret))))
	{
		// apiV1.GET("/counts", routes.Count)
		apiV1.GET("/epochs", routes.Epochs)
		apiV1.GET("/site_types", routes.SiteTypes)
		apiV1.GET("/cultures", routes.Cultures)
		apiV1.GET("/cities", routes.Cities)
		apiV1.GET("/organizations", routes.Organizations)

		apiV1.GET("/sites", routes.Sites)
		// apiV1.GET("/researches", research.Plural)
		// apiV1.GET("/authors", author.Plural)
		// apiV1.GET("/reports", report.Plural)
		// apiV1.GET("/heritages", heritage.Plural)
		// apiV1.GET("/excavations", excavation.Plural)
		// apiV1.GET("/radiocarbons", radiocarbon.Plural)
		// apiV1.GET("/artifacts", artifact.Plural)
		// apiV1.GET("/publications", publication.Plural)

		apiV1.GET("/site/:id", routes.SingleSite)
		apiV1.GET("/site/:id/researches", routes.SiteResearches)
	}

	r.Run(":8181")
}
