package routes

import (
	"api/controllers"
	"api/models"
	"api/services"

	"github.com/gin-gonic/gin"
)

// SetTickersRoutes configures the tickers routes
func SetTickersRoutes(router *gin.RouterGroup, config *models.ServerConfig) {
	stocksController := controllers.NewTickersController(services.NewTickerService(config.DB, config.Cache))

	router.GET("/", stocksController.GetTickersRecommendations)
	router.GET("/:id/logo", stocksController.GetTickerLogo)
	router.GET("/:id/overview", stocksController.GetTickerOverview)
}
