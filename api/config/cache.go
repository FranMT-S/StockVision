package config

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	Addr     string
}

var cacheConfig CacheConfig

func Cache() CacheConfig {
	if cacheConfig.Host != "" {
		return cacheConfig
	}

	host := getEnvWithDefault("REDIS_HOST", "localhost")
	port := getEnvWithDefault("REDIS_PORT", "6379")
	password := getEnvWithDefault("REDIS_PASSWORD", "")

	cacheConfig = CacheConfig{
		Host:     host,
		Port:     port,
		Password: password,
		Addr:     host + ":" + port,
	}

	return cacheConfig
}
