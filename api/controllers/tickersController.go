package controllers

import (
	apilogger "api/logger"
	"api/models"
	"api/models/filters"
	"api/models/responses"
	"api/services"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TickersController handles stock-related operations
type TickersController struct {
	tickerService services.TickerService
}

// NewTickersController creates a new tickerController
func NewTickersController(tickerService services.TickerService) TickersController {
	return TickersController{
		tickerService: tickerService,
	}
}

// GetTickerLogo retrieves the logo url of a ticker
// Path param: id (string)
func (c *TickersController) GetTickerLogo(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, models.NewResponseError("Ticker ID is required"))
		return
	}

	logo, err := c.tickerService.GetLogoUrl(ctx.Request.Context(), id)
	if err != nil {
		apilogger.Logger().Error().Err(err).Msg("Failed to retrieve ticker logo")
		ctx.JSON(http.StatusInternalServerError, models.NewResponseError("Failed to retrieve ticker logo"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": logo,
	})
}

// GetTickers retrieves a paginated list of tickers
// Query params: page (int), pageSize (int), order (asc/desc)
func (c *TickersController) GetTickersRecommendations(ctx *gin.Context) {
	filter := parseFilters(ctx)

	tickers, err := c.tickerService.GetRecommendations(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewResponseError("Failed to retrieve tickers"))
		return
	}

	if len(tickers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":  []responses.RecomendationResponse{},
			"count": 0,
		})
		return
	}

	recomendations := make([]responses.RecomendationResponse, len(tickers))

	for i, ticker := range tickers {
		recomendations[i] = responses.RecomendationResponse{
			Recomendation: ticker,
			Ticker:        ticker.Ticker,
		}
	}

	var wg sync.WaitGroup
	for i, r := range recomendations {
		wg.Add(1)
		fmt.Println(r.Ticker.ID)
		go func() {
			defer wg.Done()
			companyData, err := c.tickerService.GetCompanyData(ctx.Request.Context(), string(r.Ticker.ID))
			if err != nil {
				apilogger.Logger().Error().Err(err).Msg("Failed to retrieve ticker with ID:" + string(r.Ticker.ID))
			}

			fmt.Println("put data in recomendation" + string(r.Ticker.ID))
			recomendations[i].CompanyData = companyData
		}()
	}

	wg.Wait()

	ctx.JSON(http.StatusOK, gin.H{
		"data":  recomendations,
		"count": len(recomendations),
	})
}

// GetTickerOverview retrieves a single ticker by ID with its recommendations
// Path param: id (string)
func (c *TickersController) GetTickerOverview(ctx *gin.Context) {
	id := ctx.Param("id")

	var tickerOverview responses.CompanyOverview

	if id == "" {
		ctx.JSON(http.StatusBadRequest, models.NewResponseError("Ticker ID is required"))
		return
	}

	ticker, err := c.tickerService.GetTickerByID(ctx.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, models.NewResponseError("Ticker not found"))
			return
		}

		apilogger.Logger().Error().Err(err).Msg("Failed to retrieve ticker with ID:" + id)
		ctx.JSON(http.StatusInternalServerError, models.NewResponseError("Failed to retrieve ticker"))
		return
	}

	tickerOverview.Recomendations = ticker.Recommendations

	var wg sync.WaitGroup
	wg.Add(3)

	// Get company data
	go func() {
		defer wg.Done()
		companyData, err := c.tickerService.GetCompanyData(ctx.Request.Context(), id)
		if err != nil {
			apilogger.Logger().Error().Err(err).Msg("Failed to retrieve company data with ID:" + id)
			tickerOverview.CompanyData = models.CompanyData{}
			return
		}
		tickerOverview.CompanyData = companyData
	}()

	// Get historical prices
	go func() {
		defer wg.Done()
		historicalPrices, err := c.tickerService.GetHistoricalPrices(ctx.Request.Context(), id, "", "")
		if err != nil {
			apilogger.Logger().Error().Err(err).Msg("Failed to retrieve historical prices with ID:" + id)
			tickerOverview.HistoricalPrices = []models.HistoricalPrice{}
			return
		}

		tickerOverview.HistoricalPrices = historicalPrices
	}()

	// Get company news
	go func() {
		defer wg.Done()
		companyNews, err := c.tickerService.GetNews(ctx.Request.Context(), id, "", "")
		if err != nil {
			apilogger.Logger().Error().Err(err).Msg("Failed to retrieve company news with ID:" + id)
			tickerOverview.CompanyNews = []models.CompanyNew{}
			return
		}
		tickerOverview.CompanyNews = companyNews
	}()
	wg.Wait()

	ctx.JSON(http.StatusOK, gin.H{
		"data": tickerOverview,
	})
}

// GetRecommendations retrieves a paginated list of recommendations
// Query params: page (int, default: 1), pageSize (int, default: 10), order (asc/desc)
func (c *TickersController) GetRecommendations(ctx *gin.Context) {
	filter := parseFilters(ctx)

	recommendations, err := c.tickerService.GetRecommendations(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewResponseError("Failed to retrieve recommendations"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  recommendations,
		"count": len(recommendations),
	})
}

// parseFilters extracts pagination and ordering parameters from query string
func parseFilters(ctx *gin.Context) filters.Filters {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "0"))
	sortStr := strings.ToLower(ctx.DefaultQuery("sort", ""))

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
	}
}
