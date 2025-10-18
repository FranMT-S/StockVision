package models

// CompanyData represents company profile and market data information
type CompanyData struct {
	Symbol           string  `json:"symbol"`
	Price            float64 `json:"price"`
	MarketCap        float64 `json:"marketCap"`
	Beta             float64 `json:"beta"`
	LastDividend     float64 `json:"lastDividend"`
	Change           float64 `json:"change"`
	ChangePercentage float64 `json:"changePercentage"`
	Volume           float64 `json:"volume"`
	AverageVolume    float64 `json:"averageVolume"`
	CompanyName      string  `json:"companyName"`
	ExchangeFullName string  `json:"exchangeFullName"`
	Exchange         string  `json:"exchange"`
	Industry         string  `json:"industry"`
	Website          string  `json:"website"`
	Sector           string  `json:"sector"`
	Country          string  `json:"country"`
	Image            string  `json:"image"`
	CEO              string  `json:"ceo"`
}
