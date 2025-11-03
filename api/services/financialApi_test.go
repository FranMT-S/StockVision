package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"api/models"
	"api/services"
	CustomClient "api/services/customClient"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeFinancialCacheExpiration(t *testing.T) {

	tcc := []struct {
		name     string
		input    services.FinancialCacheExpiration
		expected services.FinancialCacheExpiration
	}{
		{
			name: "default values",
			input: services.FinancialCacheExpiration{
				HistoricalPrices: 0,
				CompanyData:      0,
			},
			expected: services.FinancialCacheExpiration{
				HistoricalPrices: 10 * time.Minute,
				CompanyData:      30 * time.Minute,
			},
		},
		{
			name: "custom values",
			input: services.FinancialCacheExpiration{
				HistoricalPrices: 5 * time.Minute,
				CompanyData:      15 * time.Minute,
			},
			expected: services.FinancialCacheExpiration{
				HistoricalPrices: 5 * time.Minute,
				CompanyData:      15 * time.Minute,
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Normalize()
			assert.Equal(t, tc.expected, result)
		})
	}

}

func TestGetCompanyDataSanatized(t *testing.T) {
	mockCompanyData := []models.CompanyData{{
		Symbol:           "AAPL",
		Price:            150.00,
		MarketCap:        2.6,
		Beta:             1.2,
		LastDividend:     0.5,
		Change:           1.5,
		ChangePercentage: 0.01,
		Volume:           1000000,
		AverageVolume:    1000000,
		CompanyName:      "Apple Inc.'; -- DROP TABLE users",
		ExchangeFullName: "NASDAQ Global Select Market\" OR 1=1 --",
		Exchange:         "NASDAQ'; EXEC xp_cmdshell('dir') --",
		Industry:         "<script>alert('XSS')</script>",
		Website:          "https://www.apple.com'; --",
		Sector:           "Technology",
		Country:          "USA\"; shutdown --",
		Image:            "https://finnhub.io/api/logo/AAPL; --",
		CEO:              "Timothy D. Cook'; DROP DATABASE test;--",
	}}

	mockCompanyDataExpected := []models.CompanyData{{
		Symbol:           "AAPL",
		Price:            150.00,
		MarketCap:        2.6,
		Beta:             1.2,
		LastDividend:     0.5,
		Change:           1.5,
		ChangePercentage: 0.01,
		Volume:           1000000,
		AverageVolume:    1000000,
		CompanyName:      "Apple Inc.&#39 ",
		ExchangeFullName: "NASDAQ Global Select Market&#34 OR 1=1 ",
		Exchange:         "NASDAQ&#39 EXEC xp_cmdshell(&#39dir&#39) ",
		Industry:         "&ltscript&gtalert(&#39XSS&#39)&lt/script&gt",
		Website:          "https://www.apple.com&#39 ",
		Sector:           "Technology",
		Country:          "USA&#34 shutdown ",
		Image:            "https://finnhub.io/api/logo/AAPL ",
		CEO:              "Timothy D. Cook&#39 DROP DATABASE test",
	}}

	mockServer := initMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stable/profile" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockCompanyData)
	})
	defer mockServer.Close()

	finghubService := &services.FinancialService{
		Client:          CustomClient.NewCustomClient(mockServer.URL),
		BaseURL:         mockServer.URL,
		Token:           "test_token",
		Cache:           nil,
		CacheExpiration: services.FinancialCacheExpiration{},
	}
	companyData, err := finghubService.GetCompanyData(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("[TestGetCompanyData] failed to get company data: %v", err)
	}

	assert.Equal(t, mockCompanyDataExpected[0], companyData)
}

func TestGetCompanyData(t *testing.T) {
	mockCompanyData := []models.CompanyData{{
		Symbol:           "AAPL",
		Price:            150.00,
		MarketCap:        2.6,
		Beta:             1.2,
		LastDividend:     0.5,
		Change:           1.5,
		ChangePercentage: 0.01,
		Volume:           1000000,
		AverageVolume:    1000000,
		CompanyName:      "Apple Inc.",
		ExchangeFullName: "NASDAQ Global Select Market",
		Exchange:         "NASDAQ",
		Industry:         "Technology",
		Website:          "https://www.apple.com",
		Sector:           "Technology",
		Country:          "USA",
		Image:            "https://finnhub.io/api/logo/AAPL",
		CEO:              "Timothy D. Cook",
	}}

	mockServer := initMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stable/profile" {
			http.NotFound(w, r)
			return
		}

		if r.URL.Query().Get("symbol") != "AAPL" {
			http.NotFound(w, r)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockCompanyData)
	})

	defer mockServer.Close()

	finghubService := &services.FinancialService{
		Client:          CustomClient.NewCustomClient(mockServer.URL),
		BaseURL:         mockServer.URL,
		Token:           "test_token",
		Cache:           nil,
		CacheExpiration: services.FinancialCacheExpiration{},
	}

	// should return error
	_, err := finghubService.GetCompanyData(context.Background(), "NVG")
	assert.ErrorContains(t, err, "[FinancialService] failed to retrieve company data id: NVG")

	//
	companyData, err := finghubService.GetCompanyData(context.Background(), "AAPL")
	assert.NoError(t, err)
	assert.Equal(t, mockCompanyData[0], companyData)
}

func TestGetHistoricalPrices(t *testing.T) {

	mockHistoricalPrices := []models.HistoricalPrice{
		{
			Symbol:  "AAPL",
			Date:    "2022-01-01",
			Open:    110.00,
			High:    120.00,
			Low:     150.00,
			Close:   160.00,
			Volume:  1000000,
			Change:  1.5,
			ChangeP: 0.01,
			Vwap:    150.00,
		},
		{
			Symbol:  "AAPL",
			Date:    "2022-01-03",
			Open:    150.00,
			High:    150.00,
			Low:     150.00,
			Close:   150.00,
			Volume:  1000000,
			Change:  1.5,
			ChangeP: 0.01,
			Vwap:    150.00,
		},

		{
			Symbol:  "NVG",
			Date:    "2022-01-02",
			Open:    110.00,
			High:    100.00,
			Low:     100.00,
			Close:   100.00,
			Volume:  1000000,
			Change:  1.0,
			ChangeP: 0.01,
			Vwap:    150.00,
		},
		{
			Symbol:  "AAPL",
			Date:    "2022-01-04",
			Open:    200.00,
			High:    200.00,
			Low:     200.00,
			Close:   200.00,
			Volume:  1000000,
			Change:  1.5,
			ChangeP: 0.01,
			Vwap:    200.00,
		},
	}

	mockServer := initMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stable/historical-price-eod/full" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		symbol := r.URL.Query().Get("symbol")
		if symbol != "AAPL" && symbol != "NVG" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]models.HistoricalPrice{})
			return
		}

		historicalPricesFiltered := make([]models.HistoricalPrice, 0)
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		for _, historicalPrice := range mockHistoricalPrices {
			if historicalPrice.Symbol == symbol {
				// filter by from and to [ >= from and <= to]
				if from != "" && to != "" {
					timeFrom, _ := time.Parse("2006-01-02", from)
					timeTo, _ := time.Parse("2006-01-02", to)
					date, _ := time.Parse("2006-01-02", historicalPrice.Date)

					if (date.After(timeFrom) || date.Equal(timeFrom)) && (date.Before(timeTo) || date.Equal(timeTo)) {
						historicalPricesFiltered = append(historicalPricesFiltered, historicalPrice)
					}
					continue
				}

				historicalPricesFiltered = append(historicalPricesFiltered, historicalPrice)
			}
		}

		json.NewEncoder(w).Encode(historicalPricesFiltered)
	})

	defer mockServer.Close()

	finhubService := &services.FinancialService{
		Client:          CustomClient.NewCustomClient(mockServer.URL),
		BaseURL:         mockServer.URL,
		Token:           "test_token",
		Cache:           nil,
		CacheExpiration: services.FinancialCacheExpiration{},
	}

	// should return error
	hp, err := finhubService.GetHistoricalPrices(context.Background(), "TEST", time.Time{}, time.Time{})
	t.Run("Should Return empty", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Empty(t, hp)
	})

	t.Run("Must be return one historical price of NVG", func(t *testing.T) {
		hp, err = finhubService.GetHistoricalPrices(context.Background(), "NVG", time.Time{}, time.Time{})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(hp))
	})

	t.Run("Must be return two historical from AAPL", func(t *testing.T) {
		hp, err = finhubService.GetHistoricalPrices(
			context.Background(), "AAPL",
			time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
			time.Date(2022, 1, 4, 0, 0, 0, 0, time.Local),
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hp))
	})

}
