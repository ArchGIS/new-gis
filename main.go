package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/ArchGIS/new-gis/assert"
	middle "github.com/ArchGIS/new-gis/middlewares"
	"github.com/ArchGIS/new-gis/routes"
	"github.com/ArchGIS/new-gis/routes/sites"
)

func init() {
	err := godotenv.Load()
	assert.Nil(err)
}

func main() {
	e := echo.New()

	e.Debug = true

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
	apiRouter.GET("/sites", sites.Plural)

	e.Logger.Fatal(e.Start(":8181"))
}
