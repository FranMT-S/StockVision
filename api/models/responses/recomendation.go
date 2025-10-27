package responses

import (
	"api/models"
)

type RecomendationResponse struct {
	Ticker      models.Ticker      `json:"ticker"`
	CompanyData models.CompanyData `json:"companyData"`
	Advice      string             `json:"advice"`
}

type CompanyOverview struct {
	CompanyData      models.CompanyData       `json:"companyData"`
	Recommendations  []models.Recommendation  `json:"recommendations"`
	HistoricalPrices []models.HistoricalPrice `json:"historicalPrices"`
	CompanyNews      []models.CompanyNew      `json:"companyNews"`
	Advice           string                   `json:"advice"`
}
