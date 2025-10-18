package config

type LogConfig struct {
	Level string
	DB    bool
}

var logConfig *LogConfig

func Log() *LogConfig {
	if logConfig == nil {
		logConfig = &LogConfig{
			Level: getEnvWithDefault("LOG_LEVEL", "info"),
			DB:    getEnvWithDefault("LOG_DB", "false") == "true",
		}
	}
	return logConfig
}
