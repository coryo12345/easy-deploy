package server

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/coryo12345/easy-deploy/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func adaptor(component templ.Component) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set("content-type", "text/html; charset=utf-8")
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
}

func (s *echoServer) RegisterServerRoutes() {
	s.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(web.StaticFiles),
	}))

	s.GET("/", adaptor(web.LoginPage()))
	s.POST("/login", s.LoginHandler)
	s.POST("/logout", s.LogoutHandler)

	authGroup := s.Group("/monitor")
	authGroup.Use(s.RequireAuth)
	authGroup.GET("", adaptor(web.MonitorPage()))
	authGroup.GET("/", adaptor(web.MonitorPage()))

}

func (s *echoServer) LoginHandler(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     X_AUTH_COOKIE,
		Value:    "password",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	})
	c.Response().Header().Set("HX-Redirect", "/monitor")
	return nil
}

func (s *echoServer) LogoutHandler(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     X_AUTH_COOKIE,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Now().AddDate(-1, 0, 0),
	})
	c.Response().Header().Set("HX-Redirect", "/")
	return nil
}
