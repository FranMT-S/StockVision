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

	advice, err := geminiai.GenerateAdvice("AAPL", getTestHistoricStock(), 20, nil)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Advice", advice)
}

func TestGeneratePredict(t *testing.T) {
	_ = godotenv.Load("../../.env")

	predictions, err := geminiai.GeneratePredict("AAPL", getTestHistoricStock(), 20, 7, nil)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Predictions", predictions)

}

func getTestHistoricStock() []models.HistoricalPrice {
	historicalData := []models.HistoricalPrice{
		{
			Symbol:  "AAPL",
			Date:    "2025-10-23",
			Open:    259.89,
			High:    260.6199,
			Low:     258.0101,
			Close:   259.58,
			Volume:  32618794,
			Change:  -0.31,
			ChangeP: -0.11928123,
			Vwap:    259.4,
		},
		{
			Symbol:  "AAPL",
			Date:    "2025-10-22",
			Open:    262.65,
			High:    262.85,
			Low:     255.43,
			Close:   258.45,
			Volume:  45015300,
			Change:  -4.2,
			ChangeP: -1.6,
			Vwap:    259.845,
		},
		{
			Symbol:  "AAPL",
			Date:    "2025-10-21",
			Open:    261.88,
			High:    265.29,
			Low:     261.83,
			Close:   262.77,
			Volume:  46695948,
			Change:  0.89,
			ChangeP: 0.33985,
			Vwap:    262.9425,
		},
		{
			Symbol:  "AAPL",
			Date:    "2025-10-20",
			Open:    255.89,
			High:    264.38,
			Low:     255.63,
			Close:   262.24,
			Volume:  90483029,
			Change:  6.36,
			ChangeP: 2.48,
			Vwap:    259.535,
		},
		{
			Symbol:  "AAPL",
			Date:    "2025-10-17",
			Open:    248.02,
			High:    253.38,
			Low:     247.27,
			Close:   252.29,
			Volume:  49147000,
			Change:  4.27,
			ChangeP: 1.72,
			Vwap:    250.24,
		},
		{
			Symbol:  "AAPL",
			Date:    "2025-10-16",
			Open:    248.25,
			High:    249.04,
			Low:     245.13,
			Close:   247.45,
			Volume:  39777000,
			Change:  -0.8,
			ChangeP: -0.32226,
			Vwap:    247.4675,
		},
		{
			Symbol:  "AAPL",
			Date:    "2025-10-15",
			Open:    249.49,
			High:    251.82,
			Low:     247.47,
			Close:   249.34,
			Volume:  33893611,
			Change:  -0.145,
			ChangeP: -0.06012265,
			Vwap:    249.53,
		},
		{
			Symbol:  "AAPL",
			Date:    "2025-10-14",
			Open:    246.6,
			High:    248.85,
			Low:     244.7,
			Close:   247.77,
			Volume:  35478000,
			Change:  1.17,
			ChangeP: 0.47445,
			Vwap:    246.98,
		},
	}

	return historicalData
}
