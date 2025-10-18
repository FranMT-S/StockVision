package models

// CompanyNew represents a news article related to a company
type CompanyNew struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Datetime int    `json:"datetime"`
	Headline string `json:"headline"`
	Image    string `json:"image"`
	Related  string `json:"related"`
	Source   string `json:"source"`
	Summary  string `json:"summary"`
	URL      string `json:"url"`
}
