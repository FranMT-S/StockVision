package services

import (
	"api/config"
	apilogger "api/logger"
	"api/models"
	CustomClient "api/services/customClient"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AnalystRatingResponse struct for analyst ratings response
type AnalystRatingResponse struct {
	Next  string                       `json:"next_page"`
	Items []models.StockRecommendation `json:"items"`
}

// AnalystRatingsServiceInterface interface for analyst ratings service
type AnalystRatingsServiceInterface interface {
	GetAll() ([]models.StockRecommendation, error)
	GetWithNext(nextPage string) (AnalystRatingResponse, error)
}

// AnalystRatingsService struct for analyst ratings service
type AnalystRatingsService struct {
	db *gorm.DB
}

// NewAnalystRatingsService to get the recommendations from the stock api
func NewAnalystRatingsService(db *gorm.DB) AnalystRatingsService {
	return AnalystRatingsService{db: db}
}

// GetAll returns all recommendations
// iterate all items until next_page is empty
func (s *AnalystRatingsService) GetAll() ([]models.StockRecommendation, error) {
	var recommendations []models.StockRecommendation
	var recommendation AnalystRatingResponse

	client := CustomClient.NewCustomClient(config.StockApi().Url)
	client.SetAuthToken(config.StockApi().Token)

	i := 1
	for {
		if err := client.Get("/list", map[string]string{"next_page": recommendation.Next}, &recommendation); err != nil {
			return nil, err
		}

		recommendations = append(recommendations, recommendation.Items...)
		if recommendation.Next == "" || strings.ToLower(recommendation.Next) == "null" {
			break
		}

		apilogger.Logger().Trace().Msg(fmt.Sprintf("Page %d Next %s", i, recommendation.Next))
		i++
	}

	return recommendations, nil
}

// GetWithNext returns the recommendations for the next page
// nextPage is the next page to fetch if its empty it will fetch the first page
func (s *AnalystRatingsService) GetWithNext(nextPage string) (AnalystRatingResponse, error) {
	var recommendation AnalystRatingResponse

	client := CustomClient.NewCustomClient(config.StockApi().Url)
	client.SetAuthToken(config.StockApi().Token)

	if err := client.Get("/list", map[string]string{"next_page": nextPage}, &recommendation); err != nil {
		return recommendation, err
	}

	for i := range recommendation.Items {
		recommendation.Items[i].Action = recommendation.Items[i].Action.Normalize()
	}

	return recommendation, nil
}

// ParseTimeNanoToRFC3339 parses a time string in RFC3339Nano format
// 2025-09-15T00:30:05.082205952Z
func ParseTimeNanoToRFC3339(t time.Time) (time.Time, error) {
	return time.Parse(time.RFC3339, t.Format(time.RFC3339))
}
