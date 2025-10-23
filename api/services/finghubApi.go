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

type FinghubService struct {
	baseURL         string
	client          CustomClient.CustomClient
	token           string
	cache           cache.ICache
	CacheExpiration FinghubCacheExpiration
}

type FinghubCacheExpiration struct {
	News time.Duration
}

func (f FinghubCacheExpiration) Normalize() FinghubCacheExpiration {
	if f.News <= 0 {
		f.News = 5 * time.Minute
	}

	return f
}

func NewFinghubService(cache cache.ICache, cacheExpiration FinghubCacheExpiration) *FinghubService {
	cacheExpiration = cacheExpiration.Normalize()

	return &FinghubService{
		client:          CustomClient.NewCustomClient(config.Finnhub().Url),
		baseURL:         config.Finnhub().Url,
		token:           config.Finnhub().Token,
		cache:           cache,
		CacheExpiration: cacheExpiration,
	}
}

func (s *FinghubService) GetNews(ctx context.Context, ticker string, from string, to string) ([]models.CompanyNew, error) {
	queryParams := map[string]string{
		"symbol": strings.ToUpper(ticker),
		"token":  s.token,
	}

	// create DateRangeFilter
	if from == "" {
		now := time.Now()
		aYearAgo := now.AddDate(-1, 0, 0)
		date := aYearAgo.Format("2006-01-02")
		queryParams["from"] = date
	}

	if to == "" {
		date := time.Now().Format("2006-01-02")
		queryParams["to"] = date
	}

	cacheKey := fmt.Sprintf("FinghubService:news:%s:%s:%s", ticker, from, to)
	expiration := s.CacheExpiration.News

	news, err := cache.GetOrLoad(ctx, s.cache, cacheKey, expiration, func() ([]models.CompanyNew, error) {
		var news []models.CompanyNew
		if err := s.client.Get("/company-news", queryParams, &news); err != nil {
			return nil, fmt.Errorf("[FinghubService] failed to retrieve news id: %s: %w", ticker, err)
		}
		return news, nil
	})

	if err != nil {
		return nil, err
	}

	for i := range news {
		news[i].Headline = sanatizer.SanatizerString(news[i].Headline).SanatizedAll().String()
		news[i].Summary = sanatizer.SanatizerString(news[i].Summary).SanatizedAll().String()
		news[i].Image = sanatizer.SanatizerString(news[i].Image).SanatizedAll().String()
		news[i].URL = sanatizer.SanatizerString(news[i].URL).SanatizedAll().String()
		news[i].Source = sanatizer.SanatizerString(news[i].Source).SanatizedAll().String()
		news[i].Related = sanatizer.SanatizerString(news[i].Related).SanatizedAll().String()
		news[i].Category = sanatizer.SanatizerString(news[i].Category).SanatizedAll().String()
		t := time.Unix(int64(news[i].Datetime), 0).UTC()
		news[i].DatetimeUTC = t.String()

	}

	return news, nil
}
