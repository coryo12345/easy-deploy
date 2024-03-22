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

	// idk why this generates random bytestrings
	// if !strings.Contains(strings.ToLower(env), "prod") {
	// 	e.Use(echoprometheus.NewMiddleware("myapp"))
	// 	e.GET("/metrics", echoprometheus.NewHandler())
	// }

}

func (s *echoServer) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(X_AUTH_COOKIE)
		if err != nil || !s.authRepo.Authenticate(cookie.Value) {
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
}
