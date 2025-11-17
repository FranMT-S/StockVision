package services

import (
	"api/cache"
	"api/database/scopes"
	"api/models"
	"api/models/filters"
	"api/models/ratings"
	"fmt"
	"sort"
	"strings"
	"time"

	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TickerService defines the interface for stock-related operations
type TickerService interface {
	GetTickers(ctx context.Context, filters filters.Filters) ([]models.Ticker, int64, error)
	GetTickerByID(ctx context.Context, id string) (*models.Ticker, error)
	GetRecommendations(ctx context.Context, filters filters.Filters) ([]models.Recommendation, error)

	// Insert operations
	InsertTickers(ctx context.Context, tickers []models.Ticker, batchSize int) (int64, error)
	InsertBrokerages(ctx context.Context, brokerages []models.Brokerage, batchSize int) (int64, error)
	InsertRecommendations(ctx context.Context, recommendations []models.Recommendation, batchSize int) (int64, error)

	// Overview operations
	GetHistoricalPrices(ctx context.Context, ticker string, from time.Time, to time.Time) ([]models.HistoricalPrice, error)
	GetLogo(ctx context.Context, ticker string) ([]byte, error)
	GetLogoUrl(ctx context.Context, ticker string) (string, error)
	GetCompanyData(ctx context.Context, ticker string) (models.CompanyData, error)
	GetNews(ctx context.Context, ticker string, from time.Time, to time.Time) ([]models.CompanyNew, error)
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
func (s *tickerService) GetTickers(ctx context.Context, filter filters.Filters) (tickers []models.Ticker, total int64, err error) {
	var cacheKey string = "tickers:total"

	query := s.db.WithContext(ctx).Model(&models.Ticker{})
	filter.Normalize()

	if filter.Query != "" {
		query = query.
			Where("id LIKE ?", strings.ToUpper(filter.Query)+"%").
			Or("company ILIKE ?", "%"+filter.Query+"%")

		cacheKey = "" // no cache
	}

	// get total items
	total, err = cache.GetOrLoad(ctx, s.cache, cacheKey, 30*time.Minute, func() (int64, error) {
		var total int64
		err = query.Count(&total).Error
		return total, err
	})

	query = query.Scopes(scopes.SortCompany(filter.Sort), scopes.Pagination(filter.Page, filter.PageSize))
	err = query.Preload("Recommendations.Brokerage").Find(&tickers).Error

	if err != nil {
		return nil, 0, fmt.Errorf("[TickerService] failed to retrieve tickers: %w", err)
	}

	// calculate sentiment for each ticker
	for i := range tickers {
		if tickers[i].Recommendations == nil {
			continue
		}
		tickerSentiment := createRatingCollection(tickers[i].Recommendations).CalculateSentiment()
		tickers[i].Sentiment = tickerSentiment.Sentiment
	}

	tickers = sortByPrefix(tickers, filter.Query)

	return tickers, total, err
}

func sortByPrefix(tickers []models.Ticker, prefix string) []models.Ticker {

	sliceCopy := make([]models.Ticker, len(tickers))
	copy(sliceCopy, tickers)

	sort.Slice(sliceCopy, func(i, j int) bool {
		firstWordFirstCompany := strings.Split(sliceCopy[i].Company, " ")[0]
		firstWordSecondCompany := strings.Split(sliceCopy[j].Company, " ")[0]

		firstCompanyHasPrefix := strings.HasPrefix(sliceCopy[i].Company, prefix)
		secondCompanyHasPrefix := strings.HasPrefix(sliceCopy[j].Company, prefix)

		if firstCompanyHasPrefix && secondCompanyHasPrefix {
			return len(firstWordFirstCompany) < len(firstWordSecondCompany)
		}

		if firstCompanyHasPrefix {
			return true
		}

		return false
	})

	return sliceCopy
}

func createRatingCollection(recommendations []models.Recommendation) ratings.RatingCollection {
	var ratingCollection ratings.RatingCollection
	for j := range recommendations {
		ratingCollection = append(ratingCollection, ratings.Rating(recommendations[j].RatingTo))
	}

	return ratingCollection
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
		return nil, fmt.Errorf("[TickerService] failed to retrieve ticker by id: %s: %w", id, err)
	}

	var ratingCollection ratings.RatingCollection

	for i := range ticker.Recommendations {
		if ticker.Recommendations[i].Brokerage.Name == "" {
			ticker.Recommendations[i].Brokerage.Name = "Anonymous"
		}

		ratingCollection = append(ratingCollection, ratings.Rating(ticker.Recommendations[i].RatingTo))
	}

	tickerSentiment := ratingCollection.CalculateSentiment()
	ticker.Sentiment = tickerSentiment.Sentiment

	return &ticker, nil
}

// InsertTickers implements TickerService interface
// InsertTickers inserts or updates tickers in the database
// If a ticker with the same ID exists, it will be updated
func (s *tickerService) InsertTickers(ctx context.Context, tickers []models.Ticker, batchSize int) (int64, error) {
	if len(tickers) == 0 {
		return 0, nil
	}

	if batchSize <= 0 {
		batchSize = 1000
	}

	db := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"company"}),
	})

	return batchFunc(db, tickers, batchSize)
}

// InsertBrokerages implements TickerService interface
// InsertBrokerages inserts or updates brokerages in the database
// If a brokerage with the same name exists, it will be updated
func (s *tickerService) InsertBrokerages(ctx context.Context, brokerages []models.Brokerage, batchSize int) (int64, error) {
	if len(brokerages) == 0 {
		return 0, nil
	}

	if batchSize <= 0 {
		batchSize = 1000
	}

	db := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	})

	return batchFunc(db, brokerages, batchSize)
}

// GetRecommendations implements TickerService interface
// GetRecommendations retrieves a paginated list of recommendations
// If pageSize and page are 0 or less, returns all recommendations
// Default pageSize is 10 and default page is 1 if not provided
func (s *tickerService) GetRecommendations(ctx context.Context, f filters.Filters) (recommendations []models.Recommendation, err error) {
	query := s.db.WithContext(ctx).Model(&models.Recommendation{}).
		Preload("Ticker").Preload("Brokerage")

	f.Normalize()

	query = query.Scopes(scopes.Pagination(f.Page, f.PageSize)).Order("ticker_id " + f.Sort.String())

	err = query.Find(&recommendations).Error
	return recommendations, err
}

// InsertRecommendations implements TickerService interface
// InsertRecommendations inserts or updates recommendations in the database
// If a recommendation with the same ID exists, it will be updated
func (s *tickerService) InsertRecommendations(ctx context.Context, recommendations []models.Recommendation, batchSize int) (int64, error) {
	if len(recommendations) == 0 {
		return 0, nil
	}

	if batchSize <= 0 {
		batchSize = 1000
	}

	for i := range recommendations {
		recommendations[i].Action = recommendations[i].Action.Normalize()
	}

	db := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ticker_id"}, {Name: "brokerage_id"}, {Name: "time"}},
		DoNothing: true,
	})

	return batchFunc(db, recommendations, batchSize)
}

// batchFunc is a generic function to execute a create operation in batches
func batchFunc[T any](db *gorm.DB, list []T, batchSize int) (int64, error) {
	var affectedRows int64

	for i := 0; i < len(list); i += batchSize {
		end := i + batchSize
		if end > len(list) {
			end = len(list)
		}

		listBatch := list[i:end]

		result := db.Create(&listBatch)

		if result.Error != nil {
			return affectedRows, result.Error
		}

		affectedRows += result.RowsAffected
	}

	return affectedRows, nil
}
