package ratings

// Rating type
type Rating string

// Sentiment types
type Sentiment string

const (
	PositiveSentiment Sentiment = "positive"
	NeutralSentiment  Sentiment = "neutral"
	NegativeSentiment Sentiment = "negative"
)

// Rating constants
const (
	RatingUnknown    Rating = ""
	Buy              Rating = "Buy"
	StrongBuy        Rating = "Strong Buy"
	Outperform       Rating = "Outperform"
	Outperformer     Rating = "Outperformer"
	Overweight       Rating = "Overweight"
	MarketOutperform Rating = "Market Outperform"
	SectorOutperform Rating = "Sector Outperform"
	SpeculativeBuy   Rating = "Speculative Buy"
	Positive         Rating = "Positive"

	Neutral       Rating = "Neutral"
	Hold          Rating = "Hold"
	EqualWeight   Rating = "Equal Weight"
	MarketPerform Rating = "Market Perform"
	SectorPerform Rating = "Sector Perform"
	InLine        Rating = "In Line"
	PeerPerform   Rating = "Peer Perform"
	SectorWeight  Rating = "Sector Weight"

	Sell               Rating = "Sell"
	Underweight        Rating = "Underweight"
	Underperform       Rating = "Underperform"
	SectorUnderperform Rating = "Sector Underperform"
	Cautious           Rating = "Cautious"
)

// RatingSentimentMap is a map of ratings to their sentiment
var RatingSentimentMap = map[Rating]Sentiment{
	RatingUnknown:    NeutralSentiment,
	Buy:              PositiveSentiment,
	StrongBuy:        PositiveSentiment,
	Outperform:       PositiveSentiment,
	Outperformer:     PositiveSentiment,
	Overweight:       PositiveSentiment,
	MarketOutperform: PositiveSentiment,
	SectorOutperform: PositiveSentiment,
	SpeculativeBuy:   PositiveSentiment,
	Positive:         PositiveSentiment,

	Neutral:       NeutralSentiment,
	Hold:          NeutralSentiment,
	EqualWeight:   NeutralSentiment,
	MarketPerform: NeutralSentiment,
	SectorPerform: NeutralSentiment,
	InLine:        NeutralSentiment,
	PeerPerform:   NeutralSentiment,
	SectorWeight:  NeutralSentiment,

	Sell:               NegativeSentiment,
	Underweight:        NegativeSentiment,
	Underperform:       NegativeSentiment,
	SectorUnderperform: NegativeSentiment,
	Cautious:           NegativeSentiment,
}
