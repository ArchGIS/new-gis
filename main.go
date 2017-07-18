package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/jmcvetta/neoism.v1"

	"github.com/ArchGIS/new-gis/assert"
	middle "github.com/ArchGIS/new-gis/middlewares"
	"github.com/ArchGIS/new-gis/routes"
)

var DB *neoism.Database

func init() {
	err := godotenv.Load()
	assert.Nil(err)
}

func main() {
	e := echo.New()

	e.Debug = true

	e.Use(middle.HandleOptions())
	e.Use(middle.AddOrigin())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", loginHandler)
	e.GET("/monuments", routes.Monuments)

	e.Logger.Fatal(e.Start(":8181"))
}
