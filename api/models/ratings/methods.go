package ratings

// RatingClassification holds the sentiment and score for a rating
type RatingClassification struct {
	Sentiment Sentiment
	Score     int // Range: -100 (very bearish) to +100 (very bullish)
}

// GetRatingClassification returns the classification for a given rating
func GetRatingClassification(r Rating) RatingClassification {
	classifications := map[Rating]RatingClassification{
		// BULLISH (Positive) - Score: 60 to 100
		StrongBuy:        {BullishSentiment, 100},
		Buy:              {BullishSentiment, 90},
		SpeculativeBuy:   {BullishSentiment, 85},
		Outperform:       {BullishSentiment, 75},
		Outperformer:     {BullishSentiment, 75},
		MarketOutperform: {BullishSentiment, 70},
		SectorOutperform: {BullishSentiment, 70},
		Overweight:       {BullishSentiment, 65},
		Positive:         {BullishSentiment, 60},

		// NEUTRAL - Score: -30 to 30
		Hold:          {NeutralSentiment, 0},
		Neutral:       {NeutralSentiment, 0},
		MarketPerform: {NeutralSentiment, 5},
		SectorPerform: {NeutralSentiment, 5},
		PeerPerform:   {NeutralSentiment, 5},
		EqualWeight:   {NeutralSentiment, 0},
		SectorWeight:  {NeutralSentiment, 0},
		InLine:        {NeutralSentiment, 0},

		// BEARISH (Negative) - Score: -100 to -40
		Underweight:  {BearishSentiment, -65},
		Underperform: {BearishSentiment, -75},
		Sell:         {BearishSentiment, -100},
	}

	if classification, exists := classifications[r]; exists {
		return classification
	}

	// Default for unknown ratings
	return RatingClassification{NeutralSentiment, 0}
}

// CalculateConsensus calculates the average score and sentiment from multiple ratings
func CalculateConsensus(ratings []Rating) (averageScore float64, sentiment Sentiment) {
	if len(ratings) == 0 {
		return 0, NeutralSentiment
	}

	total := 0
	for _, r := range ratings {
		classification := GetRatingClassification(r)
		total += classification.Score
	}

	avg := float64(total) / float64(len(ratings))

	// Determine sentiment based on average
	if avg >= 40 {
		return avg, BullishSentiment
	} else if avg <= -40 {
		return avg, BearishSentiment
	}
	return avg, NeutralSentiment
}

// GetSentimentDistribution returns the count of each sentiment type
func GetSentimentDistribution(ratings []Rating) map[Sentiment]int {
	distribution := map[Sentiment]int{
		BullishSentiment: 0,
		NeutralSentiment: 0,
		BearishSentiment: 0,
	}

	for _, r := range ratings {
		classification := GetRatingClassification(r)
		distribution[classification.Sentiment]++
	}

	return distribution
}
