package models

import (
	"api/cache"

	"gorm.io/gorm"
)

type ServerConfig struct {
	DB    *gorm.DB
	Cache cache.ICache
	Port  string
}

func NewServerConfig(db *gorm.DB, port string, cache cache.ICache) ServerConfig {
	return ServerConfig{
		DB:    db,
		Port:  port,
		Cache: cache,
	}
}
