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
		e.Response().Header().Set("content-type", "text/html; charset=utf-8")
		return component.Render(e.Request().Context(), e.Response().Writer)
	}
}

func (s *echoServer) RegisterServerRoutes() {
	s.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(web.StaticFiles),
	}))

	s.GET("/", adaptor(web.LoginPage()))
	s.POST("/login", s.LoginHandler)

	authGroup := s.Group("/monitor")
	authGroup.Use(s.RequireAuth)
	authGroup.GET("/", adaptor(web.MonitorPage()))

}

func (s *echoServer) LoginHandler(e echo.Context) error {
	return nil
}
