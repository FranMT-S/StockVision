package services_test

import (
	"testing"
	"time"

	"api/services"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeFinancialCacheExpiration(t *testing.T) {

	tcc := []struct {
		name     string
		input    services.FinancialCacheExpiration
		expected services.FinancialCacheExpiration
	}{
		{
			name: "default values",
			input: services.FinancialCacheExpiration{
				HistoricalPrices: 0,
				CompanyData:      0,
			},
			expected: services.FinancialCacheExpiration{
				HistoricalPrices: 10 * time.Minute,
				CompanyData:      30 * time.Minute,
			},
		},
		{
			name: "custom values",
			input: services.FinancialCacheExpiration{
				HistoricalPrices: 5 * time.Minute,
				CompanyData:      15 * time.Minute,
			},
			expected: services.FinancialCacheExpiration{
				HistoricalPrices: 5 * time.Minute,
				CompanyData:      15 * time.Minute,
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Normalize()
			assert.Equal(t, tc.expected, result)
		})
	}

}
