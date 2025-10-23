package models

import "time"

// Recommendation represents a stock recommendation made by a brokerage
type Recommendation struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TickerID    string    `gorm:"not null;uniqueIndex:idx_recommendation_unique;type:varchar(5)" json:"ticker_id"`
	Ticker      *Ticker   `gorm:"foreignKey:TickerID;references:ID " json:"ticker,omitempty"`
	BrokerageID uint      `gorm:"not null;uniqueIndex:idx_recommendation_unique" json:"brokerage_id"`
	Brokerage   Brokerage `gorm:"foreignKey:BrokerageID;references:ID" json:"brokerage"`
	TargetFrom  float64   `gorm:"type:decimal" json:"target_from"`
	TargetTo    float64   `gorm:"type:decimal" json:"target_to"`
	Action      Action    `gorm:"index:idx_recommendation_action" json:"action"`
	RatingFrom  string    `json:"rating_from"`
	RatingTo    string    `json:"rating_to"`
	Time        time.Time `gorm:"not null;uniqueIndex:idx_recommendation_unique" json:"time"`
}
