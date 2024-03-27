package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coryo12345/easy-deploy/web"
	"github.com/coryo12345/easy-deploy/web/components"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
	authGroup.POST("/deploy/:id", s.DeployContainerHandler)
	authGroup.POST("/refresh", s.RefreshConfig)
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
		return errorMessage(c, "Unable to authenticate. Check you have the correct password.")
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
	statuses, err := s.dockerRepo.GetStatuses(s.configRepo.GetAllServices())
	if err != nil {
		return adaptor(web.ErrorPage("Something went wrong..."))(c)
	}
	return adaptor(web.MonitorPage(statuses))(c)
}

func (s *echoServer) RefreshConfig(c echo.Context) error {
	err := s.configRepo.Refresh()
	if err != nil {
		return errorMessage(c, "Something went wrong, and we were unable to refresh the config file. You may need to restart your instance of easydeploy or verify the config file is correct.")
	}
	statuses, err := s.dockerRepo.GetStatuses(s.configRepo.GetAllServices())
	if err != nil {
		return errorMessage(c, "Something went wrong, and we were unable to reload the current services. You may need to restart your instance of easydeploy or verify the config file is correct.")
	}

	return adaptor(web.MonitorItems(statuses))(c)
}

func (s *echoServer) DeployContainerHandler(c echo.Context) error {
	id := c.Param("id")
	config, err := s.configRepo.FindEntryById(id)
	if err != nil {
		return errorMessage(c, "Unable to load config for this service.")
	}

	logs := strings.Builder{}

	defer s.dockerRepo.CleanWorkDir(config)

	errResponse := func(msg string) error {
		log.Println(err.Error())
		status, err2 := s.dockerRepo.GetStatus(config)
		if err2 != nil {
			log.Println(err2.Error())
			return adaptor(components.GlobalError("Unable to retrieve status after failure.\n" + msg))(c)
		}
		l := logs.String()
		return contentWithError(c, web.MonitorItem(status, &l), msg)
	}

	err = s.dockerRepo.CloneRepo(config, &logs)
	if err != nil {
		return errResponse("Unable to clone repository for this service")
	}

	err = s.dockerRepo.BuildImage(config, &logs)
	if err != nil {
		return errResponse("Unable to build image for this service")
	}

	err = s.dockerRepo.StopContainer(config, &logs)
	if err != nil {
		return errResponse("Unable to stop previous container for this service")
	}

	err = s.dockerRepo.DeleteContainer(config, &logs)
	if err != nil {
		return errResponse("Unable to delete previous container for this service")
	}

	err = s.dockerRepo.StartContainer(config, &logs)
	if err != nil {
		return errResponse("Unable to start container for this service")
	}

	status, err := s.dockerRepo.GetStatus(config)
	if err != nil {
		return adaptor(components.GlobalError("Unable to retrieve status after deployment"))(c)
	}
	l := logs.String()
	return adaptor(web.MonitorItem(status, &l))(c)
}
