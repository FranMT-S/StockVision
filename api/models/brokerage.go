package models

// Brokerage represents a financial brokerage firm
type Brokerage struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	Name            string           `json:"name" gorm:"type:varchar(200);unique"`
	Recommendations []Recommendation `gorm:"foreignKey:BrokerageID;references:ID" json:"recommendations,omitempty"`
}

// TableName specifies the table name for Brokerage
func (Brokerage) TableName() string {
	return "brokerages"
}
