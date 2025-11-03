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
	BaseURL         string
	Client          CustomClient.CustomClient
	Token           string
	Cache           cache.ICache
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
		Client:          CustomClient.NewCustomClient(config.Finnhub().Url),
		BaseURL:         config.Finnhub().Url,
		Token:           config.Finnhub().Token,
		Cache:           cache,
		CacheExpiration: cacheExpiration,
	}
}

func (s *FinghubService) GetNews(ctx context.Context, ticker string, from time.Time, to time.Time) ([]models.CompanyNew, error) {
	queryParams := map[string]string{
		"symbol": strings.ToUpper(ticker),
		"token":  s.Token,
	}

	var fromString string
	var toString string

	fromString = from.Format("2006-01-02")
	toString = to.Format("2006-01-02")

	// create DateRangeFilter
	if from.IsZero() {
		now := time.Now()
		aYearAgo := now.AddDate(-1, 0, 0)
		fromString = aYearAgo.Format("2006-01-02")
	}

	if to.IsZero() {
		toString = time.Now().Format("2006-01-02")
	}

	queryParams["from"] = fromString
	queryParams["to"] = toString

	cacheKey := fmt.Sprintf("FinghubService:news:%s:%s:%s", ticker, fromString, toString)
	expiration := s.CacheExpiration.News
	news, err := cache.GetOrLoad(ctx, s.Cache, cacheKey, expiration, func() ([]models.CompanyNew, error) {
		var news []models.CompanyNew
		if err := s.Client.Get("/company-news", queryParams, &news); err != nil {
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
