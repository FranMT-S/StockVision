package routes

import (
	"api/controllers"
	"api/models"
	"api/services"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router *chi.Mux, config *models.ServerConfig) {
	// Initialize services
	tickerService := services.NewTickerService(config.DB, config.Cache)

	// Initialize controllers
	tickersController := controllers.NewTickersController(tickerService, config.Cache)

	// API v1 routes
	router.Route("/api/v1", func(r chi.Router) {
		// Tickers routes
		r.Get("/tickers", tickersController.ListTickers)
		r.Get("/tickers/{id}/historical", tickersController.GetTickerHistoricalPrices)
		r.Get("/tickers/{id}/logo", tickersController.GetTickerLogo)
		r.Get("/tickers/{id}/overview", tickersController.GetTickerOverview)
		r.Get("/tickers/{id}/predictions", tickersController.GetTickerPredictions)
	})
}
