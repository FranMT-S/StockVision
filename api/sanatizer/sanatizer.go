package sanatizer

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

type SanatizerString string

func (s SanatizerString) SanatizeHTML() SanatizerString {
	return SanatizerString(SanitizeForHTML(s.String()))
}

func (s SanatizerString) SanatizeSQL() SanatizerString {
	return SanatizerString(SanitizeForSQL(s.String()))
}

func (s SanatizerString) SanatizedAll() SanatizerString {
	ss := s.SanatizeHTML().
		SanatizeSQL()

	return ss
}

func (s SanatizerString) WithMaxLength(maxLen int) SanatizerString {
	ss := s.String()
	if maxLen > 0 && utf8.RuneCountInString(ss) > maxLen {
		runes := []rune(ss)
		ss = string(runes[:maxLen])
	}

	return SanatizerString(ss)
}

func (s SanatizerString) SanatizedForLLM(maxLen int) SanatizerString {

	if maxLen <= 0 {
		maxLen = 4096
	}

	ss := strings.ReplaceAll(s.String(), "\n", " ")
	ss = strings.ReplaceAll(ss, "\r", " ")
	ss = strings.ReplaceAll(ss, "\t", " ")

	// remove keywords to change role
	re := regexp.MustCompile(`(?i)(ignore|disregard|forget all previous)`)
	ss = re.ReplaceAllString(ss, "[FILTERED]")

	// limit length
	if len(ss) > maxLen {
		ss = ss[:maxLen]
	}

	return SanatizerString(ss)
}

func (s SanatizerString) String() string {
	return string(s)
}

// SanitizeForHTML escapes special characters to prevent XSS when rendering in HTML
func SanitizeForHTML(s string) string {
	// Escapes <, >, &, ' and "
	return html.EscapeString(s)
}

// SanitizeForSQL performs a basic cleanup
// This function:
// - removes SQL comments (-- ... and /* ... */)
// - removes null characters
// - removes semicolons
// - doubles single quotes
func SanitizeForSQL(s string) string {

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\x00", "")

	// remove sql comments
	reLineComment := regexp.MustCompile(`(?i)--[^\n\r]*`)
	s = reLineComment.ReplaceAllString(s, "")

	// remove sql block comments
	reBlockComment := regexp.MustCompile(`(?is)/\*.*?\*/`)
	s = reBlockComment.ReplaceAllString(s, "")

	// remove semicolons
	s = strings.ReplaceAll(s, ";", "")

	// double single quotes
	s = strings.ReplaceAll(s, "'", "''")

	// replace multiple spaces with a single space
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")
	return s
}
