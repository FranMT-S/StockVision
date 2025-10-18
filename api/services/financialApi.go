package services

import (
	"api/cache"
	"api/config"
	"api/models"
	CustomClient "api/services/customClient"
	"context"
	"fmt"
	"strings"
	"time"
)

type FinancialCacheExpiration struct {
	HistoricalPrices time.Duration
	CompanyData      time.Duration
}

func (f FinancialCacheExpiration) Normalize() FinancialCacheExpiration {
	if f.HistoricalPrices <= 0 {
		f.HistoricalPrices = 10 * time.Minute
	}

	if f.CompanyData <= 0 {
		f.CompanyData = 30 * time.Minute
	}

	return f
}

type FinancialService struct {
	baseURL         string
	token           string
	client          CustomClient.CustomClient
	cache           cache.ICache
	CacheExpiration FinancialCacheExpiration
}

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
		return historicalPrices, nil
	})

	if err != nil {
		return nil, err
	}

	return historicalPrices, nil
}

func (s *FinancialService) GetLogo(ctx context.Context, ticker string) ([]byte, error) {
	url := fmt.Sprintf("/image-stock/%s.png", strings.ToUpper(ticker))
	logo, err := s.client.GetRaw(url, nil)
	if err != nil {
		return nil, fmt.Errorf("[FinancialService] failed to retrieve logo id: %s: %w", ticker, err)
	}

	return logo, nil
}

func (s *FinancialService) GetLogoUrl(ctx context.Context, ticker string) (string, error) {
	url := fmt.Sprintf("%s/image-stock/%s.png", s.baseURL, strings.ToUpper(ticker))

	return url, nil
}

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

	return companyData, nil
}
