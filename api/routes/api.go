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
	onboardingController := controllers.NewOnboardingController(services.NewOnboardingService(config.DB))
	// API v1 routes
	router.Route("/api/v1", func(r chi.Router) {
		// Tickers routes

		r.Route("/tickers", func(r chi.Router) {
			r.Get("/", tickersController.ListTickers)
			r.Get("/{id}/historical", tickersController.GetTickerHistoricalPrices)
			r.Get("/{id}/logo", tickersController.GetTickerLogo)
			r.Get("/{id}/overview", tickersController.GetTickerOverview)
			r.Get("/{id}/predictions", tickersController.GetTickerPredictions)
		})

		// Onboarding routes
		r.Route("/onboarding", func(r chi.Router) {
			r.Get("/", onboardingController.GetOnboarding)
			r.Patch("/", onboardingController.UpdateOnboarding)
		})
	})

}
