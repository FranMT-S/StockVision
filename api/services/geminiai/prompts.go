package geminiai

import (
	"api/models"
	"api/sanatizer"
	"fmt"
	"sort"
	"strings"
	"time"
)

// buildPrompt builds the prompt with historical data
// symbol must be a valid stock symbol max 24 characters
// action is that the model must be do "example analyze the stock for the next 7 days" max 1024 characters
// dataToAnalyze is the data to analyze max 8192 characters
// instructions are the instructions for the model max 4096 characters
// additionalInstructions are additional instructions for the model max 1024 characters
func buildPrompt(symbol string, action string, dataToAnalyze string, instructions string, additionalInstructions string) string {

	tickerSanatized := sanatizer.SanatizerString(symbol).SanatizedForLLM(24).String()
	actionSanatized := sanatizer.SanatizerString(action).SanatizedForLLM(1024).String()
	dataToAnalyzeSanatized := sanatizer.SanatizerString(dataToAnalyze).SanatizedForLLM(8192).String()
	instructionsSanatized := sanatizer.SanatizerString(instructions).SanatizedForLLM(4096).String()
	additionalInstructionsSanatized := sanatizer.SanatizerString(additionalInstructions).SanatizedForLLM(1024).String()

	prompt := fmt.Sprintf(`
	You are an expert financial analyst. Analyze the following historical data for the stock %s 
	and %s.

	DATA TO ANALYZE:
	---- START USER DATA ----
	%s
	---- END USER DATA ----

	INSTRUCTIONS:
	%s

	ADDITIONAL INSTRUCTIONS: %s.
`, tickerSanatized, actionSanatized, dataToAnalyzeSanatized, instructionsSanatized, additionalInstructionsSanatized)

	return prompt
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
	4. Predictions must have realistic variation (not static values).
	5. Volume must vary realistically.
	6. change = close - open.
	7. changePercent = (change / open) * 100.
	
	
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

	historicalStr.WriteString(formatHistoricalData(historicalData[start:]))

	return historicalStr.String()
}
