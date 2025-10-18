package responses

import (
	"api/models"
)

type RecomendationResponse struct {
	Recomendation models.Recommendation `json:"recomendation"`
	Ticker        models.Ticker         `json:"ticker"`
	CompanyData   models.CompanyData    `json:"companyData"`
}

type CompanyOverview struct {
	CompanyData      models.CompanyData       `json:"companyData"`
	Recomendations   []models.Recommendation  `json:"recomendations"`
	HistoricalPrices []models.HistoricalPrice `json:"historicalPrices"`
	CompanyNews      []models.CompanyNew      `json:"companyNews"`
}
