package controllers

import (
	"api/cache"
	apilogger "api/logger"
	"api/models"
	"api/models/filters"
	"api/models/responses"
	"api/services"
	"api/services/geminiai"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// TickersController handles stock-related operations
type TickersController struct {
	tickerService services.TickerService
	cache         cache.ICache
}

// NewTickersController creates a new tickerController
func NewTickersController(tickerService services.TickerService, cache cache.ICache) TickersController {
	return TickersController{
		tickerService: tickerService,
		cache:         cache,
	}
}

// GetTickerLogo retrieves the logo url of a ticker
// Path param: id (string)
func (c *TickersController) GetTickerLogo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		respondError(w, http.StatusBadRequest, "Ticker ID is required")
		return
	}

	logo, err := c.tickerService.GetLogoUrl(r.Context(), id)
	if err != nil {
		apilogger.Logger().Error().Err(err).Msg("Failed to retrieve ticker logo")
		respondError(w, http.StatusInternalServerError, "Failed to retrieve ticker logo")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": logo,
	})
}

// ListTickers retrieves a paginated list of tickers, with company data and recommendations
// Query params: page (int), pageSize (int), order (asc/desc)
func (c *TickersController) ListTickers(w http.ResponseWriter, r *http.Request) {
	filter := parseFilters(r)

	ctxCancel, cancelManual := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancelManual()

	tickers, total, err := c.tickerService.GetTickers(ctxCancel, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve tickers")
		return
	}

	if len(tickers) == 0 {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"data":  []responses.RecomendationResponse{},
			"total": total,
		})
		return
	}

	recomendations := make([]responses.RecomendationResponse, len(tickers))

	for i, ticker := range tickers {
		recomendations[i] = responses.RecomendationResponse{
			Ticker: ticker,
		}
	}

	var twentyDaysAgo time.Time = time.Now().AddDate(0, 0, -20)
	var wg sync.WaitGroup
	// Get company data
	for i, rec := range recomendations {
		wg.Add(1)

		go func(index int, r responses.RecomendationResponse) {
			defer wg.Done()
			companyData, err := c.tickerService.GetCompanyData(ctxCancel, string(r.Ticker.ID))
			if err != nil {
				apilogger.Logger().Error().Err(err).Msg("[ListTickers] Failed to retrieve company data with ID:" + string(r.Ticker.ID))
			}

			historicalPrices, err := c.tickerService.GetHistoricalPrices(ctxCancel, string(r.Ticker.ID), twentyDaysAgo, time.Time{})
			if err != nil {
				apilogger.Logger().Error().Err(err).Msg("[ListTickers] Failed to retrieve historical prices with ID:" + string(r.Ticker.ID))
			}

			advice, err := geminiai.GenerateAdvice(string(r.Ticker.ID), historicalPrices, 20, c.cache)
			if err != nil {
				apilogger.Logger().Error().Err(err).Msg("[ListTickers] Failed to retrieve stock analysis with ID:" + string(r.Ticker.ID))
			}

			recomendations[index].CompanyData = companyData
			recomendations[index].Advice = advice
		}(i, rec)
	}

	wg.Wait()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":  recomendations,
		"total": total,
	})
}

// GetTickerOverview retrieves a single ticker by ID with its recommendations
// Path param: id (string)
func (c *TickersController) GetTickerOverview(w http.ResponseWriter, r *http.Request) {
	var tickerOverview responses.CompanyOverview
	id := chi.URLParam(r, "id")
	from, _, err := parseDateRange(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctxCancel, cancelManual := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancelManual()

	if id == "" {
		respondError(w, http.StatusBadRequest, "Ticker ID is required")
		return
	}

	ticker, err := c.tickerService.GetTickerByID(ctxCancel, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			respondError(w, http.StatusNotFound, "Ticker not found")
			return
		}

		apilogger.Logger().Error().Err(err).Msg("[GetTickerOverview] Failed to retrieve ticker with ID:" + id)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve ticker")
		return
	}

	tickerOverview.Recommendations = ticker.Recommendations

	var wg sync.WaitGroup
	// Get company data
	wg.Add(1)
	go func() {
		defer wg.Done()
		companyData, err := c.tickerService.GetCompanyData(ctxCancel, id)
		if err != nil {
			apilogger.Logger().Error().Err(err).Msg("[GetTickerOverview] Failed to retrieve company data with ID:" + id)
			tickerOverview.CompanyData = models.CompanyData{}
			return
		}
		tickerOverview.CompanyData = companyData
	}()

	// Get historical prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		historicalPrices, err := c.tickerService.GetHistoricalPrices(ctxCancel, id, from, time.Time{})
		if err != nil {
			apilogger.Logger().Error().Err(err).Msg("[GetTickerOverview] Failed to retrieve historical prices with ID:" + id)
			tickerOverview.HistoricalPrices = []models.HistoricalPrice{}
			return
		}

		var advice string

		advice, err = geminiai.GenerateAdvice(id, historicalPrices, 20, c.cache)
		if err != nil {
			apilogger.Logger().Error().Err(err).Msg("[GetTickerOverview] Failed to retrieve stock analysis with ID:" + id)
			advice = ""
		}

		tickerOverview.Advice = advice
		tickerOverview.HistoricalPrices = historicalPrices
	}()

	// Get company news
	wg.Add(1)
	go func() {
		defer wg.Done()
		yearAgo := time.Now().AddDate(-1, 0, 0)
		companyNews, err := c.tickerService.GetNews(ctxCancel, id, yearAgo, time.Now())
		if err != nil {
			apilogger.Logger().Error().Err(err).Msg("[GetTickerOverview] Failed to retrieve company news with ID:" + id)
			tickerOverview.CompanyNews = []models.CompanyNew{}
			return
		}
		tickerOverview.CompanyNews = companyNews
	}()
	wg.Wait()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": tickerOverview,
	})
}

// GetTickerPredictions retrieves 7 days of predictions for a ticker
// Path param: id (string)
func (c *TickersController) GetTickerPredictions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var predicts []models.HistoricalPrice = make([]models.HistoricalPrice, 0)
	ctxCancel, cancelManual := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancelManual()

	if id == "" {
		respondError(w, http.StatusBadRequest, "Ticker ID is required")
		return
	}

	_, err := c.tickerService.GetTickerByID(ctxCancel, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "Ticker not found")
			return
		}

		apilogger.Logger().Error().Err(err).Msg("[GetTickerPredictions] Failed to retrieve ticker with ID:" + id)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve ticker")
		return
	}

	now := time.Now()
	before := time.Now().AddDate(0, 0, -30)

	historicalPrices, err := c.tickerService.GetHistoricalPrices(ctxCancel, id, before, now)
	if err != nil {
		apilogger.Logger().Error().Err(err).Msg("[GetTickerPredictions] Failed to retrieve historical prices with ID:" + id)
		respondError(w, http.StatusInternalServerError, "Failed in generate predictions, try again later")
		return
	}

	predicts, err = geminiai.GeneratePredict(id, historicalPrices, 20, 7, c.cache)
	if err != nil {
		apilogger.Logger().Error().Err(err).Msg("[GetTickerPredictions] Failed to retrieve stock analysis with ID:" + id)
		respondError(w, http.StatusInternalServerError, "Failed in generate predictions, try again later")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": predicts,
	})
}

// GetRecommendations retrieves a paginated list of recommendations
// Query params: page (int, default: 1), pageSize (int, default: 10), order (asc/desc)
func (c *TickersController) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	filter := parseFilters(r)

	ctxCancel, cancelManual := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancelManual()

	recommendations, err := c.tickerService.GetRecommendations(ctxCancel, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve recommendations")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":  recommendations,
		"count": len(recommendations),
	})
}

// Get The the historical prices of a ticker
func (c *TickersController) GetTickerHistoricalPrices(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	from, to, err := parseDateRange(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctxCancel, cancelManual := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancelManual()

	_, err = c.tickerService.GetTickerByID(ctxCancel, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "Ticker not found")
			return
		}

		apilogger.Logger().Error().Err(err).Msg("[GetTickerHistoricalPrices] Failed to retrieve ticker with ID:" + id)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve ticker")
		return
	}

	historicalPrices, err := c.tickerService.GetHistoricalPrices(ctxCancel, id, from, to)
	if err != nil {
		apilogger.Logger().Error().Err(err).Msg("[GetTickerHistoricalPrices] Failed to retrieve historical prices with ID:" + id)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve historical prices")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": historicalPrices,
	})
}

// parseFilters extracts pagination and ordering parameters from query string
func parseFilters(r *http.Request) filters.Filters {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	size, _ := strconv.Atoi(query.Get("size"))
	sortStr := strings.ToLower(query.Get("sort"))
	queryStr := query.Get("q")

	var sort filters.Sort
	switch sortStr {
	case "asc":
		sort = filters.ASC
	case "desc":
		sort = filters.DESC
	default:
		sort = filters.ASC
	}

	return filters.Filters{
		Page:     page,
		PageSize: size,
		Sort:     sort,
		Query:    queryStr,
	}
}

// parseDateRange extracts date range parameters from query string
// validate the format is correct
// validate that to is not before from
func parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	var err error

	var fromTime time.Time
	if from != "" {
		fromTime, err = time.Parse("2006-01-02", from)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date format: the format must be YYYY-MM-DD")
		}
	}

	var toTime time.Time
	if to != "" {
		toTime, err = time.Parse("2006-01-02", to)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to date format: the format must be YYYY-MM-DD")
		}

		if fromTime.After(toTime) {
			return time.Time{}, time.Time{}, fmt.Errorf("'To' date cannot be earlier than 'From' date")
		}
	}

	return fromTime, toTime, nil
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.NewResponseError(message))
}
