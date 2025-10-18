package config

// config store utils functions and init .env

import (
	apilogger "api/logger"
	"os"

	"strings"

	"github.com/joho/godotenv"
)

func init() {
	enviroment := strings.ToLower(os.Getenv("ENV"))
	if enviroment != "production" && enviroment != "prod" {
		_ = godotenv.Load(".env")
	}

	if enviroment == "production" || enviroment == "prod" {
		apilogger.SetLogLevel("info")
	}
}

// getEnvWithDefault gets an environment variable or returns a default value
func getEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
