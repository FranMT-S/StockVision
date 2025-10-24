package geminiai

import (
	"api/sanatizer"
	"fmt"
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
