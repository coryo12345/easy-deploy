package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/coryo12345/easy-deploy/internal/config"
	"github.com/coryo12345/easy-deploy/internal/server"
	_ "github.com/joho/godotenv/autoload"
)

type environmentVariables struct {
	configFile string
	port       int
	env        string
}

func main() {
	// read env variables
	envVars := readEnvVars()

	// initialize repositories
	configRepo, err := config.New(envVars.configFile)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	// start server
	webServer := server.New(configRepo)
	webServer.RegisterServerGlobalMiddleware(envVars.env)
	webServer.RegisterServerRoutes()
	webServer.StartServer(fmt.Sprintf(":%d", envVars.port))
}

func readEnvVars() environmentVariables {
	configFile := os.Getenv("DEPLOY_CONFIG_FILE")
	if configFile == "" {
		log.Panicf("No config file path found. DEPLOY_CONFIG_FILE is not set")
	}

	portStr := os.Getenv("DEPLOY_WEB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Panicf("DEPLOY_WEB_PORT must be defined and be an integer")
	}

	env := os.Getenv("DEPLOY_ENV_ENVIRONMENT")
	if !slices.Contains([]string{"local", "dev", "test", "prod"}, env) {
		env = "prod"
	}

	return environmentVariables{
		configFile,
		port,
		env,
	}
}
