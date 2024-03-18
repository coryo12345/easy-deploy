package server

import (
	"github.com/labstack/echo/v4/middleware"
)

func (e *echoServer) RegisterServerGlobalMiddleware(env string) {
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())

	// if !strings.Contains(strings.ToLower(env), "prod") {
	// 	e.Use(echoprometheus.NewMiddleware("myapp"))
	// 	e.GET("/metrics", echoprometheus.NewHandler())
	// }

}
