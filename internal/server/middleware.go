package server

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (e *echoServer) RegisterServerGlobalMiddleware(env string) {
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			staticIndex := strings.Index(path, "/static/")
			return staticIndex >= 0 && staticIndex < 1 || path == "/favicon.ico"
		},
		Format: "{\"method\":\"${method}\", \"uri\":\"${uri}\", \"status\":\"${status}\"}\n",
	}))

	// idk why this generates random bytestrings
	// if !strings.Contains(strings.ToLower(env), "prod") {
	// 	e.Use(echoprometheus.NewMiddleware("myapp"))
	// 	e.GET("/metrics", echoprometheus.NewHandler())
	// }

}

func (e *echoServer) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("auth middleware ran")
		return next(c)
	}
}
