package config

import (
	"os"
)

type ServerConfig struct {
	Port       string
	ClientHost string
	Env        string
}

var serverConfig *ServerConfig

func Server() *ServerConfig {
	clientHost := os.Getenv("CLIENT_HOST")

	if clientHost == "" {
		clientHost = "http://localhost:5173"
	}

	if serverConfig == nil {
		serverConfig = &ServerConfig{
			Port:       getEnvWithDefault("API_PORT", "8080"),
			ClientHost: clientHost,
			Env:        getEnvWithDefault("ENV", "development"),
		}
	}

	return serverConfig
}
