package server

import (
	"github.com/a-h/templ"
	"github.com/coryo12345/easy-deploy/web/components"
	"github.com/labstack/echo/v4"
)

func adaptor(component templ.Component) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set("content-type", "text/html; charset=utf-8")
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
}

func errorMessage(c echo.Context, msg string) error {
	return adaptor(components.GlobalError(msg))(c)
}

func contentWithError(c echo.Context, content templ.Component, errorMsg string) error {
	return adaptor(components.Fragment(content, components.GlobalError(errorMsg)))(c)
}
