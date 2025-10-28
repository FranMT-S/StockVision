package geminiai

import (
	"api/models"
	"api/sanatizer"
	"fmt"
	"strings"
)

func formatHistoricalData(data []models.HistoricalPrice) string {
	var sb strings.Builder
	for _, d := range data {
		sb.WriteString(fmt.Sprintf("%s: O:%.2f C:%.2f H:%.2f L:%.2f V:%f\n",
			sanatizer.SanatizerString(d.Date).SanatizedForLLM(10).String(),
			d.Open,
			d.Close,
			d.High,
			d.Low,
			d.Volume,
		))
	}

	return sb.String()
}
