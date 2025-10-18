package routes

import (
	"api/models"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the routes
func SetupRoutes(router *gin.Engine, config *models.ServerConfig) {
	api := router.Group("/api")

	v1 := api.Group("v1")
	tickers := v1.Group("tickers")
	SetTickersRoutes(tickers, config)

}
