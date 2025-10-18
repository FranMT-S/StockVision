package config

import (
	"os"
	"strings"
)

type ServerConfig struct {
	Port       string
	ClientHost []string
	Env        string
}

var serverConfig *ServerConfig

func Server() *ServerConfig {
	clientHost := strings.Split(os.Getenv("CLIENTS_HOST"), ",")

	if len(clientHost) == 1 && clientHost[0] == "" {
		clientHost = []string{"http://localhost:5173"}
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
