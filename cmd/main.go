package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coryo12345/easy-deploy/internal/config"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// read env variables
	configFile := os.Getenv("DEPLOY_CONFIG_FILE")
	if configFile == "" {
		log.Panicf("No config file path found. DEPLOY_CONFIG_FILE is not set")
	}

	// initialize repositories
	configRepo, err := config.New(configFile)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	fmt.Printf("%v\n", configRepo)
}
