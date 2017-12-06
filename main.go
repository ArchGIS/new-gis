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
	err := routes.InitEnv(os.Getenv("NEO4J_BOLT"))
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(middle.AddOrigin())
	r.Use(middle.HandleOptions())

	r.Any("/login", loginHandler)

	api := r.Group("/api")
	// apiV1.Use(middleware.JWT([]byte(os.Getenv(authSecret))))
	{
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

		// apiV1.GET("/site/:id", routes.SingleSite)
		// apiV1.GET("/site/:id/researches", routes.SiteResearches)
		// apiV1.GET("/site/:id/reports", routes.SiteReports)
		// apiV1.GET("/site/:id/excavations", routes.SiteExcavations)
		// apiV1.GET("/site/:id/artifacts", routes.SiteArtifacts)
		// apiV1.GET("/site/:id/radiocarbon", routes.SiteRadioCarbon)
		// apiV1.GET("/site/:id/photos", routes.SitePhotos)
		// apiV1.GET("/site/:id/topos", routes.SiteTopoplans)
	}

	r.Run(":8181")
}
