package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (e *echoServer) RegisterServerRoutes() {
	e.GET("/", func(c echo.Context) error {
		configs := e.configRepo.GetAllServices()
		c.JSON(http.StatusOK, configs)
		return nil
	})
}
