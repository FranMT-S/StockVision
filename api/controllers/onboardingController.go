package controllers

import (
	apilogger "api/logger"
	"api/models"
	"api/services"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// controller for onboarding
type OnboardingController struct {
	onboardingService services.OnboardingService
}

func NewOnboardingController(onboardingService services.OnboardingService) *OnboardingController {
	return &OnboardingController{
		onboardingService: onboardingService,
	}
}

func (c *OnboardingController) GetOnboarding(w http.ResponseWriter, r *http.Request) {

	ctxCancel, cancelManual := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancelManual()

	onboarding, err := c.onboardingService.GetOnboarding(ctxCancel)
	if err != nil {
		apilogger.Logger().Error().Err(err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve onboarding")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": onboarding,
	})
}

func (c *OnboardingController) UpdateOnboarding(w http.ResponseWriter, r *http.Request) {
	ctxCancel, cancelManual := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancelManual()

	var onboarding models.Onboarding
	err := json.NewDecoder(r.Body).Decode(&onboarding)
	if err != nil {
		apilogger.Logger().Error().Err(err)
		respondError(w, http.StatusBadRequest, "The body of the request is not valid")
		return
	}

	updatedOnboarding, err := c.onboardingService.UpdateOnboarding(ctxCancel, onboarding.OverviewStep, onboarding.OverviewDone)
	if err != nil {
		apilogger.Logger().Error().Err(err)
		respondError(w, http.StatusInternalServerError, "Failed to update onboarding")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": updatedOnboarding,
	})
}
