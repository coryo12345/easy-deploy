package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const X_AUTH_COOKIE = "x-easy-deploy-auth"

func (s *echoServer) RegisterServerGlobalMiddleware() {
	s.Use(middleware.Recover())
	s.Use(middleware.Secure())
	s.Use(middleware.Gzip())
	s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			staticIndex := strings.Index(path, "/static/")
			return staticIndex >= 0 && staticIndex < 1 || path == "/favicon.ico"
		},
		Format: "{\"method\":\"${method}\", \"uri\":\"${uri}\", \"status\":\"${status}\"}\n",
	}))
}

func (s *echoServer) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(X_AUTH_COOKIE)
		if err != nil {
			return c.Redirect(http.StatusFound, "/")
		}
		token := cookie.Value
		valid := s.jwtBuilder.VerifyToken(token)
		if valid {
			return next(c)
		} else {
			return c.Redirect(http.StatusFound, "/")
		}
	}
}
