package config

import (
	"strings"
)

type DatabaseConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	SSLMode    string
	SchemaName string
}

var databaseConfig *DatabaseConfig

// Database returns the database configuration
func Database() *DatabaseConfig {

	sslMode := "disable"
	sslModeEnv := getEnvWithDefault("DB_SSL", "false")
	if strings.ToLower(sslModeEnv) == "true" {
		sslMode = "require"
	}

	if databaseConfig == nil {
		databaseConfig = &DatabaseConfig{
			Host:       getEnvWithDefault("DB_HOST", "localhost"),
			Port:       getEnvWithDefault("DB_PORT", "26257"),
			User:       getEnvWithDefault("DB_USER", "root"),
			Password:   getEnvWithDefault("DB_PASSWORD", ""),
			DBName:     getEnvWithDefault("DB_NAME", "stocksvision"),
			SSLMode:    sslMode,
			SchemaName: getEnvWithDefault("DB_SCHEMA", "public"),
		}
	}

	return databaseConfig
}
