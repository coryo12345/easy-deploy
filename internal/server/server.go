package server

import (
	"github.com/coryo12345/easy-deploy/internal/config"
	"github.com/labstack/echo/v4"
)

type Server interface {
	RegisterServerRoutes()
	RegisterServerGlobalMiddleware(env string)
	StartServer(host string)
}

type echoServer struct {
	*echo.Echo
	configRepo config.ConfigRepository
}

func New(configRepo config.ConfigRepository) Server {
	e := echo.New()

	return &echoServer{
		Echo:       e,
		configRepo: configRepo,
	}
}

func (s *echoServer) StartServer(host string) {
	s.Logger.Fatal(s.Start(host))
}
