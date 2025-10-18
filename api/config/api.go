package config

type StockApiConfig struct {
	Url   string
	Token string
}

var stockApiConfigInstance *StockApiConfig

// GetStockApiConfig returns the stockApiConfig instance
func StockApi() *StockApiConfig {
	if stockApiConfigInstance == nil {
		stockApiConfigInstance = &StockApiConfig{
			Url:   getEnvWithDefault("STOCK_API_URL", ""),
			Token: getEnvWithDefault("STOCK_API_TOKEN", ""),
		}
	}

	return stockApiConfigInstance
}

type FinancialModelingConfig struct {
	Url   string
	Token string
}

var financialModelingConfigInstance *FinancialModelingConfig

// GetFinancialModelingConfig returns the financialModelingConfig instance
func FinancialModeling() *FinancialModelingConfig {
	if financialModelingConfigInstance == nil {
		financialModelingConfigInstance = &FinancialModelingConfig{
			Url:   getEnvWithDefault("FINANCIAL_BASE_URL", ""),
			Token: getEnvWithDefault("FINANCIAL_TOKEN", ""),
		}
	}

	return financialModelingConfigInstance
}

type FinnhubConfig struct {
	Url   string
	Token string
}

var finnhubConfigInstance *FinnhubConfig

// GetFinnhubConfig returns the finnhubConfig instance
func Finnhub() *FinnhubConfig {
	if finnhubConfigInstance == nil {
		finnhubConfigInstance = &FinnhubConfig{
			Url:   getEnvWithDefault("FINHUB_BASE_URL", ""),
			Token: getEnvWithDefault("FINHUB_TOKEN", ""),
		}
	}

	return finnhubConfigInstance
}
