package filters

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateFloat(t *testing.T) {

	tcc := []struct {
		input    float64
		expected float64
		decimals int
	}{
		{
			input:    1.2345,
			expected: 1.23,
			decimals: 2,
		},
		{
			input:    2.2345545,
			expected: 2.234,
			decimals: 3,
		},
		{
			input:    0.245,
			expected: 0.0,
			decimals: 0,
		},
	}

	for _, tC := range tcc {
		t.Run(fmt.Sprintf("TruncateFloat(%f, %d)", tC.input, tC.decimals), func(t *testing.T) {
			assert.Equal(t, tC.expected, TruncateFloat(tC.input, tC.decimals))
		})
	}

}
