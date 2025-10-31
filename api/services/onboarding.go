package services

import (
	"api/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type OnboardingService struct {
	db *gorm.DB
}

func NewOnboardingService(db *gorm.DB) OnboardingService {
	return OnboardingService{
		db: db,
	}
}

// get the onboarding data
//
// in this case we only have a one data
func (s *OnboardingService) GetOnboarding(ctx context.Context) (*models.Onboarding, error) {
	var onboarding models.Onboarding
	err := s.db.WithContext(ctx).FirstOrCreate(&onboarding, models.Onboarding{ID: 1}).Error
	if err != nil {
		return nil, fmt.Errorf("[OnboardingService] failed to get the onboarding: %w", err)
	}

	return &onboarding, nil
}

// update the onboarding data
//
// in this case we only have a one data to update
func (s *OnboardingService) UpdateOnboarding(ctx context.Context, step int, done bool) (*models.Onboarding, error) {

	var onboarding models.Onboarding
	err := s.db.WithContext(ctx).FirstOrCreate(&onboarding, models.Onboarding{ID: 1}).Error
	if err != nil {
		return nil, fmt.Errorf("[OnboardingService] The onboarding with ID %d was not found: %w", 1, err)
	}

	onboarding.OverviewStep = step
	onboarding.OverviewDone = done
	err = s.db.WithContext(ctx).Save(&onboarding).Error
	if err != nil {
		return nil, fmt.Errorf("[OnboardingService] failed to update onboarding: %w", err)
	}

	return &onboarding, nil
}
