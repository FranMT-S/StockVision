package geminiai_test

import (
	"api/models"
	"api/services/geminiai"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestGenerateAdviceAndStocks(t *testing.T) {
	_ = godotenv.Load("../../.env")

	prediction, err := geminiai.GenerateAdviceAndStocksPredict("AAPL", []models.HistoricalPrice{})
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Advice", prediction.Advice)
	fmt.Println("Stocks", prediction.StocksNextWeek)
}
