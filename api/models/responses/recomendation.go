package responses

import (
	"api/models"
)

type RecomendationResponse struct {
	Ticker      models.Ticker      `json:"ticker"`
	CompanyData models.CompanyData `json:"companyData"`
}

type CompanyOverview struct {
	CompanyData      models.CompanyData       `json:"companyData"`
	Recommendations  []models.Recommendation  `json:"recommendations"`
	HistoricalPrices []models.HistoricalPrice `json:"historicalPrices"`
	CompanyNews      []models.CompanyNew      `json:"companyNews"`
}

type PredictResponse struct {
	Advice         string                   `json:"advice"`
	StocksNextWeek []models.HistoricalPrice `json:"stocksNextWeek"`
}
