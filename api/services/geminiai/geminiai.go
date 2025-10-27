package geminiai

import (
	"api/cache"
	"api/config"
	"api/models"
	"api/sanatizer"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
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
		temp := float32(0.2) // low for consistency
		topP := float32(0.7)
		topK := float32(30)
		thinkingBudget := int32(0)

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
	generationConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   predictSchema,
	}

	key := fmt.Sprintf("GeminiAI:predict:%s-%s", symbol, time.Now().Format("2006-01-02"))
	expiration := 30 * time.Minute

	result, err := cache.GetOrLoad(ctx, c, key, expiration, func() ([]models.HistoricalPrice, error) {
		if daysToPredict > 14 {
			daysToPredict = 14
		}

		if daysToPredict < 1 {
			daysToPredict = 7
		}

		if daysToAnalyze > 30 {
			daysToAnalyze = 30
		}

		if daysToAnalyze < 1 {
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

		var prediction struct {
			StocksNextWeek []models.HistoricalPrice `json:"stocksNextWeek"`
		}

		err = json.Unmarshal([]byte(result.Text()), &prediction)
		if err != nil {
			return make([]models.HistoricalPrice, 0), fmt.Errorf("[GeminiAI] cannot unmarshal JSON: %s", result.Text())
		}

		return prediction.StocksNextWeek, nil
	})

	if err != nil {
		return make([]models.HistoricalPrice, 0), err
	}

	return result, nil
}

// buildPromptAdvice builds the prompt for advice
// symbol must be a valid stock symbol max 24 characters
// dayToAnalyze is the number of days to analyze max 30 days
func buildPromptAdvice(symbol string, historicalData []models.HistoricalPrice, dayToAnalyze int) string {
	action := "generate advice for the stock"
	data := buildHistoricalDataString(symbol, historicalData, dayToAnalyze)
	instructions := `	
	1. Analyze recent price trends, volatility, and trading volume.
	2. Determine if the current market [BEHAVIOR] suggests BUY, HOLD, or SELL.
	3. Provide a short, clear, and realistic recommendation (max 10 words).
	4. Include a one-sentence [JUSTIFICATION] for the advice.
	5. The advice must be in English, plain text (no markdown or HTML).
	6. Do NOT restate or summarize the data.
	7. The format must be [Behavior]. [JUSTIFICATION]

	Respond ONLY with the advice and justification.`

	additionalInstructions := "Provide only the advice. Don't mention the labels"

	prompt := buildPrompt(symbol, action, data, instructions, additionalInstructions)
	return prompt
}

// buildPrompt builds the prompt with historical data
// date format must be 2006-01-02
func buildPredictPromp(symbol string, historicalData []models.HistoricalPrice, dayToAnalyze int, dayToPredict int) string {
	action := fmt.Sprintf("generate predictions for the next %d trading days", dayToPredict)
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	instructions := fmt.Sprintf(`
	1. Analyze price trends, volume, and volatility.
	2. Consider technical patterns (support, resistance, moving averages).
	3. Generate realistic predictions based on historical data.
	4. Provide a clear recommendation: BUY, HOLD, or SELL with justification.
	5. Predictions must have realistic variation (not static values).
	6. Volume must vary realistically.
	7. Calculate 'change' and 'changePercent' correctly based on the previous day.
	
	Respond with the predictions for %d trading days, the prediction must be start from the date %s.
	`, dayToPredict, tomorrow)

	data := buildHistoricalDataString(symbol, historicalData, dayToAnalyze)
	additionalInstructions := "Provides predictions. Don't mention the labels."

	prompt := buildPrompt(symbol, action, data, instructions, additionalInstructions)

	return prompt
}

// buildHistoricalDataString builds the historical data string
// date format must be 2006-01-02
// maxDays is the max number of days to analyze, the limit is 30
func buildHistoricalDataString(symbol string, historicalData []models.HistoricalPrice, maxDays int) string {
	var historicalStr strings.Builder
	historicalStr.WriteString(fmt.Sprintf("Historical data of %s:\n", symbol))

	start := 0

	if len(historicalData) == 0 {
		return ""
	}

	if maxDays < 1 {
		maxDays = 7
	}

	if maxDays > 30 {
		maxDays = 30
	}

	sort.Slice(historicalData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", historicalData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", historicalData[j].Date)

		return dateI.After(dateJ)
	})

	if len(historicalData) > maxDays {
		start = len(historicalData) - maxDays
	}

	for i := start; i < len(historicalData); i++ {
		h := historicalData[i]
		historicalStr.WriteString(fmt.Sprintf(
			"- %s: Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %.0f, Change: %.2f%%\n",
			sanatizer.SanatizerString(h.Date).SanatizedForLLM(10).String(), h.Open, h.High, h.Low, h.Close, h.Volume, h.ChangeP,
		))
	}

	return historicalStr.String()
}
