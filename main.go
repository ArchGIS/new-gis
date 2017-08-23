package main

import (
	"os"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/ArchGIS/new-gis/assert"
	middle "github.com/ArchGIS/new-gis/middlewares"
	"github.com/ArchGIS/new-gis/routes"
)

func init() {
	err := godotenv.Load()
	assert.Nil(err)
}

func main() {
	err := routes.InitEnv(os.Getenv("Neo4jHost"))
	assert.Nil(err)
	e := echo.New()

	e.Debug = true

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middle.AddOrigin())
	e.Use(middle.HandleOptions())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Match(
		[]string{"POST", "GET"},
		"/login",
		loginHandler,
	)

	apiRouter := e.Group("/api")
	apiRouter.Use(middleware.JWT([]byte(os.Getenv(authSecret))))

	apiRouter.GET("/counts", routes.Count)
	apiRouter.GET("/epochs", routes.Epochs)
	apiRouter.GET("/site_types", routes.SiteTypes)
	apiRouter.GET("/cultures", routes.Cultures)
	apiRouter.GET("/cities", routes.Cities)
	apiRouter.GET("/organizations", routes.Organizations)

	apiRouter.GET("/sites", routes.Sites)
	// apiRouter.GET("/researches", research.Plural)
	// apiRouter.GET("/authors", author.Plural)
	// apiRouter.GET("/reports", report.Plural)
	// apiRouter.GET("/heritages", heritage.Plural)
	// apiRouter.GET("/excavations", excavation.Plural)
	// apiRouter.GET("/radiocarbons", radiocarbon.Plural)
	// apiRouter.GET("/artifacts", artifact.Plural)
	// apiRouter.GET("/publications", publication.Plural)

	apiRouter.GET("/site/:id", routes.SingleSite)

	e.Logger.Fatal(e.Start(":8181"))
}
