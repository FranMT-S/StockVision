package geminiai

import "google.golang.org/genai"

type StockPredict struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
	Vwap   float64 `json:"vwap"`
}

var predictSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"stocksNextWeek": {
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"date":   {Type: genai.TypeString, Description: "Trading date (YYYY-MM-DD)"},
					"open":   {Type: genai.TypeNumber, Description: "Predicted closing price"},
					"close":  {Type: genai.TypeNumber, Description: "Predicted closing price"},
					"high":   {Type: genai.TypeNumber, Description: "Predicted high price"},
					"low":    {Type: genai.TypeNumber, Description: "Predicted low price"},
					"volume": {Type: genai.TypeNumber, Description: "Predicted trading volume"},
					"vwap":   {Type: genai.TypeNumber, Description: "Predicted volume-weighted average price"},
				},
			},
		},
	},
}
