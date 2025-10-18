package services

import (
	"api/models"
	"context"
)

type HistoricalPriceService interface {
	GetHistoricalPrices(ctx context.Context, ticker string, from string, to string) ([]models.HistoricalPrice, error)
}

type LogoService interface {
	GetLogo(ctx context.Context, ticker string) ([]byte, error)
	GetLogoUrl(ctx context.Context, ticker string) (string, error)
}

type CompanyDataService interface {
	GetCompanyData(ctx context.Context, ticker string) (models.CompanyData, error)
}

type CompanyNewsService interface {
	GetNews(ctx context.Context, ticker string, from string, to string) ([]models.CompanyNew, error)
}
