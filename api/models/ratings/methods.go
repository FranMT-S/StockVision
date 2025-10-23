package ratings

// RatingClassification holds the sentiment and score for a rating
type SentimentScore struct {
	Sentiment     Sentiment `json:"sentiment"`
	PositiveCount int       `json:"positive_count"`
	NeutralCount  int       `json:"neutral_count"`
	NegativeCount int       `json:"negative_count"`
	TotalCount    int       `json:"total_count"`
	PositiveRatio float64   `json:"positive_ratio"`
	NegativeRatio float64   `json:"negative_ratio"`
	Score         float64   `json:"score"` // -1 to 1, where 1 is most positive
}

// RatingCollection is a slice of ratings
type RatingCollection []Rating

// CalculateSentiment analyzes an array of ratings and returns the overall sentiment
func (rc RatingCollection) CalculateSentiment() SentimentScore {
	var positiveCount, neutralCount, negativeCount int

	// Count each sentiment type
	for _, rating := range rc {
		sentiment, exists := RatingSentimentMap[rating]
		if !exists {
			// Unknown ratings are treated as neutral
			neutralCount++
			continue
		}

		switch sentiment {
		case PositiveSentiment:
			positiveCount++
		case NeutralSentiment:
			neutralCount++
		case NegativeSentiment:
			negativeCount++
		default:
			neutralCount++
		}
	}

	totalCount := len(rc)
	if totalCount == 0 {
		return SentimentScore{
			Sentiment: NeutralSentiment,
		}
	}

	// Calculate ratios
	positiveRatio := float64(positiveCount) / float64(totalCount)
	negativeRatio := float64(negativeCount) / float64(totalCount)

	// Calculate score (-1 to 1)
	// Score = (positive - negative) / total
	score := (float64(positiveCount) - float64(negativeCount)) / float64(totalCount)

	// Determine overall sentiment
	var overallSentiment Sentiment
	if score > 0.2 {
		overallSentiment = PositiveSentiment
	} else if score < -0.2 {
		overallSentiment = NegativeSentiment
	} else {
		overallSentiment = NeutralSentiment
	}

	return SentimentScore{
		Sentiment:     overallSentiment,
		PositiveCount: positiveCount,
		NeutralCount:  neutralCount,
		NegativeCount: negativeCount,
		TotalCount:    totalCount,
		PositiveRatio: positiveRatio,
		NegativeRatio: negativeRatio,
		Score:         score,
	}
}

// GetRatingSentiment returns the sentiment for a single rating
func GetRatingSentiment(rating Rating) Sentiment {
	sentiment, exists := RatingSentimentMap[rating]
	if !exists {
		return NeutralSentiment
	}
	return sentiment
}

// CalculateWeightedSentiment calculates sentiment with time-based weights
// More recent ratings have higher weight
func (rc RatingCollection) CalculateWeightedSentiment(weights []float64) SentimentScore {
	if len(rc) != len(weights) {
		// Fallback to unweighted if weights don't match
		return rc.CalculateSentiment()
	}

	var positiveScore, neutralScore, negativeScore, totalWeight float64

	for i, rating := range rc {
		weight := weights[i]
		sentiment, exists := RatingSentimentMap[rating]
		if !exists {
			neutralScore += weight
			totalWeight += weight
			continue
		}

		switch sentiment {
		case PositiveSentiment:
			positiveScore += weight
		case NeutralSentiment:
			neutralScore += weight
		case NegativeSentiment:
			negativeScore += weight
		}
		totalWeight += weight
	}

	if totalWeight == 0 {
		return SentimentScore{
			Sentiment: NeutralSentiment,
		}
	}

	// Calculate weighted score
	score := (positiveScore - negativeScore) / totalWeight

	// Determine overall sentiment
	var overallSentiment Sentiment
	if score > 0.2 {
		overallSentiment = PositiveSentiment
	} else if score < -0.2 {
		overallSentiment = NegativeSentiment
	} else {
		overallSentiment = NeutralSentiment
	}

	// Calculate counts for reporting
	positiveCount := int(positiveScore / totalWeight * float64(len(rc)))
	negativeCount := int(negativeScore / totalWeight * float64(len(rc)))
	neutralCount := len(rc) - positiveCount - negativeCount

	return SentimentScore{
		Sentiment:     overallSentiment,
		PositiveCount: positiveCount,
		NeutralCount:  neutralCount,
		NegativeCount: negativeCount,
		TotalCount:    len(rc),
		PositiveRatio: positiveScore / totalWeight,
		NegativeRatio: negativeScore / totalWeight,
		Score:         score,
	}
}
