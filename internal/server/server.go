package server

import (
	"github.com/coryo12345/easy-deploy/internal/auth"
	"github.com/coryo12345/easy-deploy/internal/config"
	"github.com/labstack/echo/v4"
)

type Server interface {
	RegisterServerRoutes()
	RegisterServerGlobalMiddleware()
	StartServer(host string)
}

type echoServer struct {
	*echo.Echo
	configRepo config.ConfigRepository
	authRepo   auth.AuthRepository
	jwtBuilder auth.JwtBuilder
}

func New(configRepo config.ConfigRepository, authRepo auth.AuthRepository, jwtBuilder auth.JwtBuilder) Server {
	e := echo.New()

	return &echoServer{
		Echo:       e,
		configRepo: configRepo,
		authRepo:   authRepo,
		jwtBuilder: jwtBuilder,
	}
}

func (s *echoServer) StartServer(host string) {
	s.Logger.Fatal(s.Start(host))
}
