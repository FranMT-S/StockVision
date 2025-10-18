package ratings

// Rating type
type Rating string

// Sentiment types
type Sentiment string

const (
	BullishSentiment Sentiment = "Bullish"
	NeutralSentiment Sentiment = "Neutral"
	BearishSentiment Sentiment = "Bearish"
)

// Rating constants
const (
	Buy              Rating = "Buy"
	Outperform       Rating = "Outperform"
	Neutral          Rating = "Neutral"
	EqualWeight      Rating = "Equal Weight"
	Underweight      Rating = "Underweight"
	Overweight       Rating = "Overweight"
	MarketPerform    Rating = "Market Perform"
	Hold             Rating = "Hold"
	Underperform     Rating = "Underperform"
	Positive         Rating = "Positive"
	Sell             Rating = "Sell"
	SectorPerform    Rating = "Sector Perform"
	MarketOutperform Rating = "Market Outperform"
	InLine           Rating = "In Line"
	Outperformer     Rating = "Outperformer"
	PeerPerform      Rating = "Peer Perform"
	StrongBuy        Rating = "Strong Buy"
	SectorOutperform Rating = "Sector Outperform"
	SectorWeight     Rating = "Sector Weight"
	SpeculativeBuy   Rating = "Speculative Buy"
)
