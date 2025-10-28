package geminiai

import (
	"api/cache"
	"api/config"
	"api/models"
	"api/models/filters"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/genai"
)

// GenerateAdviceAndStocksPredict generates predictions and advice using Gemini AI
// the stock predict is for the next 7 days
//
//	with a limit of 30 to analyze
func GenerateAdvice(symbol string, historicalData []models.HistoricalPrice, daysToAnalyze int, c cache.ICache) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.GeminiAi().Token,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", err
	}

	if len(historicalData) == 0 {
		return "UNKNOWN. We don't have enough data to generate advice", nil
	}

	key := fmt.Sprintf("GeminiAI:advice:%s-%s", symbol, time.Now().Format("2006-01-02"))
	expiration := 10 * time.Minute

	result, err := cache.GetOrLoad(ctx, c, key, expiration, func() (string, error) {
		if daysToAnalyze > 30 {
			daysToAnalyze = 30
		}

		if daysToAnalyze < 1 {
			daysToAnalyze = 7
		}

		// Config Model
		temp := float32(0.2)
		topP := float32(0.7)
		topK := float32(30)
		thinkingBudget := int32(30)

		config := &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(
				"You are a quantitative analyst. Provide objective analysis based on data, without speculation.",
				genai.RoleUser,
			),
			Temperature:     &temp,
			TopP:            &topP,
			TopK:            &topK,
			MaxOutputTokens: 512, // short advice
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &thinkingBudget, // answers fast
			},
		}

		result, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(buildPromptAdvice(symbol, historicalData, daysToAnalyze)),
			config,
		)

		if err != nil {
			return "", fmt.Errorf("error: %w", err)
		}

		advice := result.Text()
		return advice, nil
	})

	if err != nil {
		return "", err
	}

	return result, nil
}

// GeneratePredict generates predictions using Gemini AI
// the stock predict is for the next 7 days
//
//	with a limit of 30 days to analyze
//	with a limit of 14 days to predict
func GeneratePredict(symbol string, historicalData []models.HistoricalPrice, daysToAnalyze int, daysToPredict int, c cache.ICache) ([]models.HistoricalPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.GeminiAi().Token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return make([]models.HistoricalPrice, 0), err
	}

	// Config Model
	temp := float32(0.2) // low for consistency
	topP := float32(0.7)
	topK := float32(30)
	thinkingBudget := int32(50)

	// Config Model
	generationConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   predictSchema,
		SystemInstruction: genai.NewContentFromText(
			"You are a quantitative analyst. Provide objective analysis based on data, without speculation.",
			genai.RoleUser,
		),
		Temperature: &temp,
		TopP:        &topP,
		TopK:        &topK,
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget, // answers fast
		},
	}

	key := fmt.Sprintf("GeminiAI:predict:%s-%s", symbol, time.Now().Format("2006-01-02"))
	expiration := 30 * time.Minute

	result, err := cache.GetOrLoad(ctx, c, key, expiration, func() ([]models.HistoricalPrice, error) {
		if daysToPredict > 14 {
			daysToPredict = 14
		}

		if daysToPredict < 7 {
			daysToPredict = 7
		}

		if daysToAnalyze > 30 {
			daysToAnalyze = 30
		}

		if daysToAnalyze < 7 {
			daysToAnalyze = 7
		}

		result, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(buildPredictPromp(symbol, historicalData[:], daysToAnalyze, daysToPredict)),
			generationConfig,
		)

		if err != nil {
			return nil, err
		}

		var historicalPredict []models.HistoricalPrice
		var prediction struct {
			StocksNextWeek []StockPredict `json:"stocksNextWeek"`
		}

		err = json.Unmarshal([]byte(result.Text()), &prediction)
		if err != nil {
			return make([]models.HistoricalPrice, 0), fmt.Errorf("[GeminiAI] cannot unmarshal JSON: %s", result.Text())
		}

		for _, p := range prediction.StocksNextWeek {
			historicalPredict = append(historicalPredict, models.HistoricalPrice{
				Symbol:  symbol,
				Date:    p.Date,
				Open:    filters.TruncateFloat(p.Open, 2),
				High:    filters.TruncateFloat(p.High, 2),
				Low:     filters.TruncateFloat(p.Low, 2),
				Close:   filters.TruncateFloat(p.Close, 2),
				Volume:  filters.TruncateFloat(p.Volume, 0),
				Change:  filters.TruncateFloat(p.Close-p.Open, 2),
				ChangeP: filters.TruncateFloat((p.Close-p.Open)/p.Open, 5),
				Vwap:    filters.TruncateFloat(p.Vwap, 4),
			})
		}

		return historicalPredict, nil
	})

	if err != nil {
		return make([]models.HistoricalPrice, 0), err
	}

	return result, nil
}
