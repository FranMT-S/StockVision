package geminiai

import "google.golang.org/genai"

var adviceSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"advice": {Type: genai.TypeString},
	},
}

var predictSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
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
}
