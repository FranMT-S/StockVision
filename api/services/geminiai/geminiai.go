package geminiai

import (
	"api/config"
	"api/models"
	"api/models/responses"
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
func GenerateAdviceAndStocksPredict(symbol string, historicalData []models.HistoricalPrice) (*responses.PredictResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.GeminiAi().Token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	// Config Model
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"advice": {Type: genai.TypeString},
				"stocksNextWeek": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"symbol":        {Type: genai.TypeString},
							"date":          {Type: genai.TypeString},
							"open":          {Type: genai.TypeNumber},
							"high":          {Type: genai.TypeNumber},
							"low":           {Type: genai.TypeNumber},
							"close":         {Type: genai.TypeNumber},
							"volume":        {Type: genai.TypeNumber},
							"change":        {Type: genai.TypeNumber},
							"changePercent": {Type: genai.TypeNumber},
							"vwap":          {Type: genai.TypeNumber},
						},
					},
				},
			},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(buildPrompt(symbol, historicalData)),
		config,
	)

	if err != nil {
		return nil, err
	}

	var prediction responses.PredictResponse

	err = json.Unmarshal([]byte(result.Text()), &prediction)
	if err != nil {
		return nil, fmt.Errorf("[GeminiAI] cannot unmarshal JSON: %s", result.Text())
	}

	return &prediction, nil
}

// buildPrompt builds the prompt with historical data
// date format must be 2006-01-02
func buildPrompt(symbol string, historicalData []models.HistoricalPrice) string {
	var historicalStr strings.Builder
	historicalStr.WriteString(fmt.Sprintf("Historical data of %s:\n", symbol))

	// take the last 15 days (or less if there are less data)
	start := 0
	maxDays := 15

	// sort desc
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

	tickerSanatized := sanatizer.SanatizerString(symbol).SanatizedForLLM(4096).String()

	prompt := fmt.Sprintf(`You are an expert financial analyst. Analyze the following historical data for the stock %s 
and generate predictions for the next 7 trading days.

DATA TO ANALYZE:
---- START USER DATA ----
%s
---- END USER DATA ----

INSTRUCTIONS:
1. Analyze price trends, volume, and volatility.
2. Consider technical patterns (support, resistance, moving averages).
3. Generate realistic predictions based on historical data.
4. Provide a clear recommendation: BUY, HOLD, or SELL with justification.
5. Predictions must have realistic variation (not static values).
6. Volume must vary realistically.
7. Calculate 'change' and 'changePercent' correctly based on the previous day.
8. The advice must be concise and clear, short max 10 words.
9. the advice must be in english, and must be a string without format html or markdown.

Respond with the predictions for 7 trading days and a substantiated investment advice.

ADDITIONAL INSTRUCTIONS: Provides prediction and advice only. Don't mention the labels.
`, tickerSanatized, historicalStr.String())

	return prompt
}
