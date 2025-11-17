package models

import (
	"time"
)

// StockRecommendation struct for stock recommendation
// Ticker is the stock ticker
type StockRecommendation struct {
	Ticker     string         `json:"ticker"`
	Company    string         `json:"company"`
	Action     Action         `json:"action"`
	Brokerage  string         `json:"brokerage"`
	TargetFrom CurrencyString `json:"target_from"`
	TargetTo   CurrencyString `json:"target_to"`
	RatingFrom string         `json:"rating_from"`
	RatingTo   string         `json:"rating_to"`
	Time       time.Time      `json:"time"`
}
