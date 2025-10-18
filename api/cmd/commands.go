package cmd

import (
	"api/database"
	apilogger "api/logger"
	"api/models"
	"api/services"
	"context"

	"github.com/spf13/cobra"
)

var fillDbCmd = &cobra.Command{
	Use:   "fill-db",
	Short: "Fill database with initial data",
	Long:  `Run fill-db to fill the database with initial data from the Stock API`,
	Run:   fillDb,
}

// FillDb fills the database with initial data from the Stock API
// run after the database is initialized
func fillDb(cmd *cobra.Command, args []string) {
	db, err := database.GetDB()
	if err != nil {
		apilogger.Logger().Err(err).Msg("failed to get database instance")
		return
	}

	analystRatingsService := services.NewAnalystRatingsService(db.DB)
	stockRecommendations, err := analystRatingsService.GetAll()
	if err != nil {
		apilogger.Logger().Err(err).Msg("failed to get recommendations")
		return
	}

	tickerService := services.NewTickerService(db.DB, nil)

	// clean and prepare the entities for insertion
	tickers, brokerages := cleanAndPrepareEntities(stockRecommendations)
	tickerService.InsertTickers(context.Background(), tickers)
	tickerService.InsertBrokerages(context.Background(), brokerages)

	// create the map of brokerages with ids
	brokeragesWithIdsMap := createBrokeragesWithIdsMap(brokerages)
	recommendations := createRecommendations(stockRecommendations, brokeragesWithIdsMap)
	tickerService.InsertRecommendations(context.Background(), recommendations)

	apilogger.Logger().Info().Msg("Database filled successfully")
}

// cleanAndPrepareEntities cleans and prepares the entities for insertion and remove the brokerages with empty name
func cleanAndPrepareEntities(stockRecommendations []services.StockRecommendation) ([]models.Ticker, []models.Brokerage) {
	var brokerages []models.Brokerage
	var tickers []models.Ticker
	var brokeragesMap = make(map[string]bool)
	var tickersMap = make(map[string]bool)

	for _, recommendation := range stockRecommendations {
		if recommendation.Ticker != "" {
			if _, ok := tickersMap[recommendation.Ticker]; !ok {
				ticker := models.Ticker{
					ID:      models.TickerID(recommendation.Ticker),
					Company: recommendation.Company,
				}

				tickersMap[recommendation.Ticker] = true
				tickers = append(tickers, ticker)
			}
		}

		if recommendation.Brokerage != "" {
			if _, ok := brokeragesMap[recommendation.Brokerage]; !ok {
				brokeragesMap[recommendation.Brokerage] = true

				brokerage := models.Brokerage{
					Name: recommendation.Brokerage,
				}
				brokerages = append(brokerages, brokerage)
			}
		}
	}

	return tickers, brokerages
}

func createBrokeragesWithIdsMap(brokerages []models.Brokerage) map[string]uint {
	var brokeragesMap = make(map[string]uint)
	for _, brokerage := range brokerages {
		if brokerage.Name != "" {
			brokeragesMap[brokerage.Name] = brokerage.ID
		}
	}

	return brokeragesMap
}

func createRecommendations(stockRecommendations []services.StockRecommendation, brokeragesWithIdsMap map[string]uint) []models.Recommendation {
	var recommendations []models.Recommendation
	for _, recommendation := range stockRecommendations {
		// // if the brokerage not exists, skip the recommendation
		// // we want to keep the recommendations with the brokerages id to prioritize confidence that theys provide
		// if brokeragesWithIdsMap[recommendation.Brokerage] == 0 {
		// 	continue
		// }

		recommendation := models.Recommendation{
			TickerID:    recommendation.Ticker,
			BrokerageID: brokeragesWithIdsMap[recommendation.Brokerage],
			TargetFrom:  recommendation.TargetFrom,
			TargetTo:    recommendation.TargetTo,
			Action:      recommendation.Action,
			RatingFrom:  recommendation.RatingFrom,
			RatingTo:    recommendation.RatingTo,
			Time:        recommendation.Time,
		}
		recommendations = append(recommendations, recommendation)
	}

	return recommendations
}
