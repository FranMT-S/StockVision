package models

// HistoricalPrice represents company historical price of stock from a company
type HistoricalPrice struct {
	Symbol  string  `json:"symbol"`
	Date    string  `json:"date"`
	Open    float64 `json:"open"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Close   float64 `json:"close"`
	Volume  float64 `json:"volume"`
	Change  float64 `json:"change"`
	ChangeP float64 `json:"changePercent"`
	Vwap    float64 `json:"vwap"`
}
