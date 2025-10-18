package services_test

import (
	"testing"
	"time"

	"api/services"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeFinghubCacheExpiration(t *testing.T) {

	tcc := []struct {
		name     string
		input    services.FinghubCacheExpiration
		expected services.FinghubCacheExpiration
	}{
		{
			name: "default values",
			input: services.FinghubCacheExpiration{
				News: 0,
			},
			expected: services.FinghubCacheExpiration{
				News: 5 * time.Minute,
			},
		},
		{
			name: "custom values",
			input: services.FinghubCacheExpiration{
				News: 10 * time.Minute,
			},
			expected: services.FinghubCacheExpiration{
				News: 10 * time.Minute,
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
