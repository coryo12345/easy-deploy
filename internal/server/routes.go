package server

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/coryo12345/easy-deploy/internal/docker"
	"github.com/coryo12345/easy-deploy/web"
	"github.com/coryo12345/easy-deploy/web/components"
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
	authGroup.GET("", s.MonitorPageHandler)
	authGroup.GET("/", s.MonitorPageHandler)

}

func (s *echoServer) LoginHandler(c echo.Context) error {
	password := c.FormValue("password")
	valid := s.authRepo.Authenticate(password)

	if valid {
		token, err := s.jwtBuilder.NewToken()
		if err != nil {
			return c.String(http.StatusInternalServerError, "something went wrong")
		}

		c.SetCookie(&http.Cookie{
			Name:     X_AUTH_COOKIE,
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
		})
		c.Response().Header().Set("HX-Redirect", "/monitor")
		return nil
	} else {
		c.Response().Header().Set("HX-Retarget", "#global-error")
		return adaptor(components.ErrorMessage("Unable to authenticate. Check you have the correct password."))(c)
	}
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

func (s *echoServer) MonitorPageHandler(c echo.Context) error {
	statuses, err := docker.GetStatuses(s.configRepo.GetAllServices())
	if err != nil {
		return adaptor(web.ErrorPage("Something went wrong..."))(c)
	}
	return adaptor(web.MonitorPage(statuses))(c)
}
