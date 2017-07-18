package middlewares

import "github.com/labstack/echo"

func AddOrigin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			return next(c)
		}
	}
}

func HandleOptions() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "OPTIONS" {
				c.Response().Header().Set("Allow", "OPTIONS, GET, POST")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
				return next(c)
			}

			return next(c)
		}
	}
}
