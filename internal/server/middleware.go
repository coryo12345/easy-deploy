package server

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (e *echoServer) RegisterServerGlobalMiddleware(env string) {
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())

	// idk why this generates random bytestrings
	// if !strings.Contains(strings.ToLower(env), "prod") {
	// 	e.Use(echoprometheus.NewMiddleware("myapp"))
	// 	e.GET("/metrics", echoprometheus.NewHandler())
	// }

}

func (e *echoServer) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("auth middleware ran")
		fmt.Println("auth middleware ran -fmt")
		return next(c)
	}
}
