package models

// Onboarding struct
//
// This struct represents the onboarding table in the database
// used to determine what process show to user
type Onboarding struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	OverviewStep int  `json:"overviewStep" gorm:"default:1"`
	OverviewDone bool `json:"overviewDone" gorm:"default:false"`
}

func (Onboarding) TableName() string {
	return "onboarding"
}
