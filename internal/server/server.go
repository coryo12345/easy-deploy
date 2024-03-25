package server

import (
	"github.com/coryo12345/easy-deploy/internal/auth"
	"github.com/coryo12345/easy-deploy/internal/config"
	"github.com/coryo12345/easy-deploy/internal/docker"
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
	dockerRepo docker.DockerRepository
}

func New(configRepo config.ConfigRepository, authRepo auth.AuthRepository, jwtBuilder auth.JwtBuilder, dockerRepo docker.DockerRepository) Server {
	e := echo.New()

	return &echoServer{
		Echo:       e,
		configRepo: configRepo,
		authRepo:   authRepo,
		jwtBuilder: jwtBuilder,
		dockerRepo: dockerRepo,
	}
}

func (s *echoServer) StartServer(host string) {
	s.Logger.Fatal(s.Start(host))
}
