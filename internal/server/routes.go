package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/coryo12345/easy-deploy/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func adaptor(component templ.Component) func(e echo.Context) error {
	return func(e echo.Context) error {
		return component.Render(e.Request().Context(), e.Response().Writer)
	}
}

func (e *echoServer) RegisterServerRoutes() {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(web.StaticFiles),
	}))

	e.GET("/", adaptor(web.Base()))

	authGroup := e.Group("/monitor")
	authGroup.Use(e.RequireAuth)
	authGroup.GET("/", adaptor(web.HomePage()))

}
