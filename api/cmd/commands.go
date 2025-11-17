package cmd

import (
	"api/database"
	apilogger "api/logger"
	"api/models"
	"api/services"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var fillDbCmd = &cobra.Command{
	Use:   "fill-db",
	Short: "Fill database with initial data, can send --json flag to get the data in json format",
	Long:  `Run fill-db to fill the database with initial data from the Stock API, can send --json flag to get the data in json format`,
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

	fmt.Println("Get recommendations")
	stockRecommendations, err := getRecommendationsData(cmd, db)
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to get recommendations")
		return err
	}

	// insert tickers and brokerages
	tickerService := services.NewTickerService(db.DB, nil)

	// clean and prepare the entities for insertion
	fmt.Println("Clean data")
	tickers, brokerages := cleanAndPrepareEntities(stockRecommendations)

	// insert tickers and brokerages
	fmt.Println("Insert tickers")
	_, err = tickerService.InsertTickers(context.Background(), tickers, 1000)
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to insert tickers")
		return err
	}
	fmt.Println("Insert tickers brokerages")
	_, err = tickerService.InsertBrokerages(context.Background(), brokerages, 1000)
	if err != nil {
		apilogger.Logger().Err(err).Msg("[fillDb] failed to insert brokerages")
		return err
	}

	// create the map of brokerages with ids
	brokeragesWithIdsMap := createBrokeragesWithIdsMap(brokerages)
	recommendations := createRecommendations(stockRecommendations, brokeragesWithIdsMap)

	fmt.Println("Insert recommendations")
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
func cleanAndPrepareEntities(stockRecommendations []models.StockRecommendation) ([]models.Ticker, []models.Brokerage) {
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

func getRecommendationsData(cmd *cobra.Command, db *database.Database) ([]models.StockRecommendation, error) {
	// get recommendations stock
	analystRatingsService := services.NewAnalystRatingsService(db.DB)
	jsonPath, _ := cmd.Flags().GetString("json")

	var stockRecommendations []models.StockRecommendation
	var err error

	var getStocksFunc func() ([]models.StockRecommendation, error) = analystRatingsService.GetAll

	if jsonPath != "" {
		getStocksFunc = func() ([]models.StockRecommendation, error) {
			return getRecommendationsFromJson(jsonPath)
		}
	}

	stockRecommendations, err = getStocksFunc()
	if err != nil {
		return nil, err
	}

	return stockRecommendations, nil
}

func getRecommendationsFromJson(path string) ([]models.StockRecommendation, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var recommendations []models.StockRecommendation
	err = json.Unmarshal(file, &recommendations)
	if err != nil {
		return nil, err
	}

	return recommendations, nil
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

func createRecommendations(stockRecommendations []models.StockRecommendation, brokeragesWithIdsMap map[string]uint) []models.Recommendation {
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
