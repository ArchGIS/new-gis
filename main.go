package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/jmcvetta/neoism.v1"

	"github.com/ArchGIS/new-gis/assert"
	middle "github.com/ArchGIS/new-gis/middlewares"
)

var DB *neoism.Database

func init() {
	err := godotenv.Load()
	assert.Nil(err)

	neoHost := os.Getenv("Neo4jHost")
	DB, err = neoism.Connect(neoHost + "db/data")
	assert.Nil(err)
}

func main() {
	e := echo.New()

	e.Debug = true

	e.Use(middle.HandleOptions())
	e.Use(middle.AddOrigin())
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Recover())

	e.POST("/login", loginHandler)

	e.Logger.Fatal(e.Start(":8181"))
}
