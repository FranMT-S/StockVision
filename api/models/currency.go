package models

import (
	"strconv"
	"strings"
)

type CurrencyString string

// CurrencyToFloat sanitizes and converts currency string to float64
func (c *CurrencyString) CurrencyToFloat() float64 {
	if *c == "" {
		return 0
	}

	// Remove dollar sign, spaces, and commas
	sanitized := strings.ReplaceAll(string(*c), "$", "")
	sanitized = strings.ReplaceAll(sanitized, ",", "")
	sanitized = strings.TrimSpace(sanitized)

	// Convert to float64
	result, err := strconv.ParseFloat(sanitized, 64)
	if err != nil {
		return 0
	}

	return result
}
