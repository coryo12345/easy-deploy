package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/coryo12345/easy-deploy/internal/auth"
	"github.com/coryo12345/easy-deploy/internal/config"
	"github.com/coryo12345/easy-deploy/internal/docker"
	"github.com/coryo12345/easy-deploy/internal/server"
	_ "github.com/joho/godotenv/autoload"
)

type environmentVariables struct {
	configFile  string
	host        string
	port        int
	env         string
	webPassword string
	workDir     string
}

func main() {
	// read env variables
	envVars := readEnvVars()

	// initialize repositories
	configRepo, err := config.New(envVars.configFile)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	authRepo := auth.NewAuthRepository(envVars.webPassword)
	jwtBuilder := auth.NewJwtBuilder(envVars.env)
	dockerRepo := docker.New(envVars.workDir)

	// TODO need to run init command from config

	// start server
	webServer := server.New(configRepo, authRepo, jwtBuilder, dockerRepo)
	webServer.RegisterServerGlobalMiddleware()
	webServer.RegisterServerRoutes()
	webServer.StartServer(fmt.Sprintf("%s:%d", envVars.host, envVars.port))
}

func readEnvVars() environmentVariables {
	configFile := os.Getenv("DEPLOY_CONFIG_FILE")
	if configFile == "" {
		log.Panic("No config file path found. DEPLOY_CONFIG_FILE is not set")
	}

	portStr := os.Getenv("DEPLOY_WEB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Panic("DEPLOY_WEB_PORT must be defined and be an integer")
	}

	env := os.Getenv("DEPLOY_ENV_ENVIRONMENT")
	if !slices.Contains([]string{"local", "dev", "test", "prod"}, env) {
		env = "prod"
	}

	host := ""
	if env == "local" {
		host = "localhost"
	}

	password := os.Getenv("DEPLOY_WEB_PASSWORD")
	if password == "" {
		log.Panic("DEPLOY_WEB_PASSWORD must be defined!")
	}

	workDir := os.Getenv("DEPLOY_WORK_DIR")
	if workDir == "" {
		log.Panic("DEPLOY_WORK_DIR must be defined")
	}

	return environmentVariables{
		configFile:  configFile,
		host:        host,
		port:        port,
		env:         env,
		webPassword: password,
		workDir:     workDir,
	}
}
