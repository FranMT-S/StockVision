package cmd

import (
	"api/database"
	apilogger "api/logger"
	"api/models"
	"api/services"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var fillDbCmd = &cobra.Command{
	Use:   "fill-db",
	Short: "Fill database with initial data",
	Long:  `Run fill-db to fill the database with initial data from the Stock API`,
	RunE:  fillDb,
}

// FillDb fills the database with initial data from the Stock API
// run after the database is initialized
func fillDb(cmd *cobra.Command, args []string) error {
	db, err := database.GetDB()
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to get database instance")
		return err
	}
	fmt.Println("Start fillDb")

	analystRatingsService := services.NewAnalystRatingsService(db.DB)
	stockRecommendations, err := analystRatingsService.GetAll()
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to get recommendations")
		return err
	}

	tickerService := services.NewTickerService(db.DB, nil)

	// clean and prepare the entities for insertion
	tickers, brokerages := cleanAndPrepareEntities(stockRecommendations)
	_, err = tickerService.InsertTickers(context.Background(), tickers, 1000)
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to insert tickers")
		return err
	}
	_, err = tickerService.InsertBrokerages(context.Background(), brokerages, 1000)
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to insert brokerages")
		return err
	}

	// create the map of brokerages with ids
	brokeragesWithIdsMap := createBrokeragesWithIdsMap(brokerages)
	recommendations := createRecommendations(stockRecommendations, brokeragesWithIdsMap)

	_, err = tickerService.InsertRecommendations(context.Background(), recommendations, 1000)
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to insert recommendations")
		return err
	}

	apilogger.Logger().Info().Msg("Database filled successfully")
	fmt.Println("Database filled successfully")
	return nil
}

// cleanAndPrepareEntities cleans and prepares the entities for insertion and remove the brokerages with empty name
func cleanAndPrepareEntities(stockRecommendations []services.StockRecommendation) ([]models.Ticker, []models.Brokerage) {
	var brokerages []models.Brokerage = make([]models.Brokerage, 0)
	var tickers []models.Ticker = make([]models.Ticker, 0)
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
	var recommendations []models.Recommendation = make([]models.Recommendation, 0)
	for _, recommendation := range stockRecommendations {

		parsedTime, err := ParseTimeNanoToRFC3339(recommendation.Time)
		if err != nil {
			apilogger.Logger().Err(err).Msg("[createRecommendations] failed to parse time, recommendation: " + recommendation.Ticker)
			continue
		}

		recommendation := models.Recommendation{
			TickerID:    recommendation.Ticker,
			BrokerageID: brokeragesWithIdsMap[recommendation.Brokerage],
			TargetFrom:  recommendation.TargetFrom.CurrencyToFloat(),
			TargetTo:    recommendation.TargetTo.CurrencyToFloat(),
			Action:      recommendation.Action,
			RatingFrom:  recommendation.RatingFrom,
			RatingTo:    recommendation.RatingTo,
			Time:        parsedTime,
		}

		recommendations = append(recommendations, recommendation)
	}

	return recommendations
}
