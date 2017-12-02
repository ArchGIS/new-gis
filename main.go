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

	// e.Validator = &CustomValidator{validator: validator.New()}

	r.Use(middle.AddOrigin())
	r.Use(middle.HandleOptions())

	r.Any("/login", loginHandler)

	apiV1 := r.Group("/api")
	// apiV1.Use(middleware.JWT([]byte(os.Getenv(authSecret))))
	{
		apiV1.POST("/graphql", routes.Graphql)

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
