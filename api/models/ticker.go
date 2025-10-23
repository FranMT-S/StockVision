package models

import (
	"api/models/ratings"
	"strings"
)

// Ticker represents a stock ticker symbol
type Ticker struct {
	ID              TickerID          `json:"id" gorm:"primaryKey;type:varchar(5)"`
	Company         string            `gorm:"not null;index:idx_ticker_company;type:varchar(200)" json:"company"`
	Recommendations []Recommendation  `gorm:"foreignKey:TickerID;references:ID" json:"recommendations,omitempty"`
	Sentiment       ratings.Sentiment `json:"sentiment" gorm:"-"`
}

type TickerID string

func (t TickerID) String() string {
	return strings.ToUpper(string(t))
}
