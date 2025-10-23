package services

import (
	"api/cache"
	"api/config"
	"api/models"
	"api/sanatizer"
	CustomClient "api/services/customClient"
	"context"
	"fmt"
	"strings"
	"time"
)

// FinancialCacheExpiration represents the cache expiration values for financial data
type FinancialCacheExpiration struct {
	HistoricalPrices time.Duration
	CompanyData      time.Duration
}

// Normalize normalizes the cache expiration values  in case of invalid values
func (f FinancialCacheExpiration) Normalize() FinancialCacheExpiration {
	if f.HistoricalPrices <= 0 {
		// f.HistoricalPrices = 10 * time.Minute
		f.HistoricalPrices = 60 * time.Minute
	}

	if f.CompanyData <= 0 {
		// f.CompanyData = 30 * time.Minute
		f.CompanyData = 60 * time.Minute
	}

	return f
}

// FinancialService represents a service for financial data
type FinancialService struct {
	baseURL         string
	token           string
	client          CustomClient.CustomClient
	cache           cache.ICache
	CacheExpiration FinancialCacheExpiration
}

// NewFinancialService creates a new FinancialService instance
func NewFinancialService(cache cache.ICache, cacheExpiration FinancialCacheExpiration) *FinancialService {
	cacheExpiration = cacheExpiration.Normalize()

	return &FinancialService{
		baseURL: config.FinancialModeling().Url,
		token:   config.FinancialModeling().Token,
		client:  CustomClient.NewCustomClient(config.FinancialModeling().Url),
		cache:   cache,
		CacheExpiration: FinancialCacheExpiration{
			HistoricalPrices: cacheExpiration.HistoricalPrices,
			CompanyData:      cacheExpiration.CompanyData,
		},
	}
}

// GetHistoricalPrices returns the historical prices of a company from symbol ticker
func (s *FinancialService) GetHistoricalPrices(ctx context.Context, ticker string, from string, to string) ([]models.HistoricalPrice, error) {
	params := map[string]string{
		"symbol": strings.ToUpper(ticker),
		"apikey": s.token,
	}

	if from != "" {
		params["from"] = from
	}

	if to != "" && from != "" {
		params["to"] = to
	}

	key := fmt.Sprintf("FinancialService:historical_prices:%s:%s:%s", ticker, from, to)
	expiration := s.CacheExpiration.HistoricalPrices

	historicalPrices, err := cache.GetOrLoad(ctx, s.cache, key, expiration, func() ([]models.HistoricalPrice, error) {
		var historicalPrices []models.HistoricalPrice
		if err := s.client.Get("/stable/historical-price-eod/full", params, &historicalPrices); err != nil {
			return nil, fmt.Errorf("[FinancialService] failed to retrieve historical prices id: %s: %w", ticker, err)
		}
		fmt.Println("Historical prices retrieved from API")
		return historicalPrices, nil
	})

	if err != nil {
		return nil, err
	}

	return historicalPrices, nil
}

// GetLogo returns the logo of a company as a byte array
func (s *FinancialService) GetLogo(ctx context.Context, ticker string) ([]byte, error) {
	url := fmt.Sprintf("/image-stock/%s.png", strings.ToUpper(ticker))
	logo, err := s.client.GetRaw(url, nil)
	if err != nil {
		return nil, fmt.Errorf("[FinancialService] failed to retrieve logo id: %s: %w", ticker, err)
	}

	return logo, nil
}

// GetLogoUrl returns the logo url of a company
func (s *FinancialService) GetLogoUrl(ctx context.Context, ticker string) (string, error) {
	url := fmt.Sprintf("%s/image-stock/%s.png", s.baseURL, strings.ToUpper(ticker))

	return url, nil
}

// GetCompanyData returns the company data of a company from symbol ticker
func (s *FinancialService) GetCompanyData(ctx context.Context, ticker string) (models.CompanyData, error) {
	params := map[string]string{
		"symbol": strings.ToUpper(ticker),
		"apikey": s.token,
	}

	key := fmt.Sprintf("FinancialService:company_data:%s", ticker)
	expiration := s.CacheExpiration.CompanyData

	companyData, err := cache.GetOrLoad(ctx, s.cache, key, expiration, func() (models.CompanyData, error) {
		var companyData []models.CompanyData
		if err := s.client.Get("/stable/profile", params, &companyData); err != nil {
			return models.CompanyData{}, fmt.Errorf("[FinancialService] failed to retrieve company data id: %s: %w", ticker, err)
		}
		return companyData[0], nil
	})

	if err != nil {
		return models.CompanyData{}, err
	}

	companyData.Website = sanatizer.SanatizerString(companyData.Website).SanatizedAll().String()
	companyData.ExchangeFullName = sanatizer.SanatizerString(companyData.ExchangeFullName).SanatizedAll().String()
	companyData.Exchange = sanatizer.SanatizerString(companyData.Exchange).SanatizedAll().String()
	companyData.CompanyName = sanatizer.SanatizerString(companyData.CompanyName).SanatizedAll().String()
	companyData.Industry = sanatizer.SanatizerString(companyData.Industry).SanatizedAll().String()
	companyData.Sector = sanatizer.SanatizerString(companyData.Sector).SanatizedAll().String()
	companyData.Country = sanatizer.SanatizerString(companyData.Country).SanatizedAll().String()
	companyData.Image = sanatizer.SanatizerString(companyData.Image).SanatizedAll().String()
	companyData.CEO = sanatizer.SanatizerString(companyData.CEO).SanatizedAll().String()

	return companyData, nil
}
