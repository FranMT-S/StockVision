package services

import (
	"api/cache"
	"api/database/scopes"
	"api/models"
	"api/models/filters"

	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TickerService defines the interface for stock-related operations
type TickerService interface {
	GetTickers(ctx context.Context, filters filters.Filters) ([]models.Ticker, error)
	GetTickerByID(ctx context.Context, id string) (*models.Ticker, error)
	GetRecommendations(ctx context.Context, filters filters.Filters) ([]models.Recommendation, error)

	// Insert operations
	InsertTickers(ctx context.Context, tickers []models.Ticker) (int64, error)
	InsertBrokerages(ctx context.Context, brokerages []models.Brokerage) (int64, error)
	InsertRecommendations(ctx context.Context, recommendations []models.Recommendation) (int64, error)

	// Overview operations
	GetHistoricalPrices(ctx context.Context, ticker string, from string, to string) ([]models.HistoricalPrice, error)
	GetLogo(ctx context.Context, ticker string) ([]byte, error)
	GetLogoUrl(ctx context.Context, ticker string) (string, error)
	GetCompanyData(ctx context.Context, ticker string) (models.CompanyData, error)
	GetNews(ctx context.Context, ticker string, from string, to string) ([]models.CompanyNew, error)
}

type tickerService struct {
	db *gorm.DB
	HistoricalPriceService
	LogoService
	CompanyDataService
	CompanyNewsService
	cache cache.ICache
}

// NewTickerService creates a new instance of TickerService
func NewTickerService(db *gorm.DB, cache cache.ICache) TickerService {
	financialApi := NewFinancialService(cache, FinancialCacheExpiration{})
	finhubApi := NewFinghubService(cache, FinghubCacheExpiration{})

	return &tickerService{
		db:                     db,
		HistoricalPriceService: financialApi,
		LogoService:            financialApi,
		CompanyDataService:     financialApi,
		CompanyNewsService:     finhubApi,
		cache:                  cache,
	}
}

// GetTickers implements TickerService interface
// GetTickers retrieves a paginated list of tickers
// If pageSize and page are 0 or less, returns all tickers
func (s *tickerService) GetTickers(ctx context.Context, f filters.Filters) (tickers []models.Ticker, err error) {
	query := s.db.WithContext(ctx).Model(&models.Ticker{})

	f.Normalize()

	query = query.Scopes(scopes.SortCompany(f.Sort), scopes.Pagination(f.Page, f.PageSize))
	err = query.Find(&tickers).Error
	return tickers, err
}

// GetTickerByID implements TickerService interface
// GetTickerByID retrieves a single ticker by ID with its recommendations preloaded
func (s *tickerService) GetTickerByID(ctx context.Context, id string) (*models.Ticker, error) {
	var ticker models.Ticker
	err := s.db.WithContext(ctx).
		Preload("Recommendations.Brokerage").
		Where("id = ?", id).
		First(&ticker).Error

	if err != nil {
		return nil, err
	}
	return &ticker, nil
}

// InsertTickers implements TickerService interface
// InsertTickers inserts or updates tickers in the database
// If a ticker with the same ID exists, it will be updated
func (s *tickerService) InsertTickers(ctx context.Context, tickers []models.Ticker) (int64, error) {
	if len(tickers) == 0 {
		return 0, nil
	}

	result := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"company"}),
	}).Create(&tickers)

	return result.RowsAffected, result.Error
}

// InsertBrokerages implements TickerService interface
// InsertBrokerages inserts or updates brokerages in the database
// If a brokerage with the same name exists, it will be updated
func (s *tickerService) InsertBrokerages(ctx context.Context, brokerages []models.Brokerage) (int64, error) {
	if len(brokerages) == 0 {
		return 0, nil
	}

	result := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&brokerages)

	return result.RowsAffected, result.Error
}

// GetRecommendations implements TickerService interface
// GetRecommendations retrieves a paginated list of recommendations
// If pageSize and page are 0 or less, returns all recommendations
// Default pageSize is 10 and default page is 1 if not provided
func (s *tickerService) GetRecommendations(ctx context.Context, f filters.Filters) (recommendations []models.Recommendation, err error) {
	query := s.db.WithContext(ctx).Model(&models.Recommendation{}).
		Preload("Ticker").Preload("Brokerage")

	f.Normalize()

	// Apply pagination

	query = query.Scopes(scopes.Pagination(f.Page, f.PageSize)).Order("ticker_id " + f.Sort.String())

	err = query.Find(&recommendations).Error
	return recommendations, err
}

// InsertRecommendations implements TickerService interface
// InsertRecommendations inserts or updates recommendations in the database
// If a recommendation with the same ID exists, it will be updated
func (s *tickerService) InsertRecommendations(ctx context.Context, recommendations []models.Recommendation) (int64, error) {
	if len(recommendations) == 0 {
		return 0, nil
	}

	for i := range recommendations {
		recommendations[i].Action = recommendations[i].Action.Normalize()
	}

	result := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ticker_id"}, {Name: "brokerage_id"}, {Name: "time"}},
		DoNothing: true,
	}).Create(&recommendations)

	return result.RowsAffected, result.Error
}
