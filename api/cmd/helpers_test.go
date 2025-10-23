package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeNano(t *testing.T) {

	tcc := []struct {
		input  string
		output string
	}{
		{
			input:  "2025-09-15T00:30:05.082205952Z",
			output: "2025-09-15 00:30:05 +0000 UTC",
		},
		{
			input:  "2025-10-06T00:30:12.961009226Z",
			output: "2025-10-06 00:30:12 +0000 UTC",
		},
		{
			input:  "2025-08-25T00:30:04.718800464Z",
			output: "2025-08-25 00:30:04 +0000 UTC",
		},
	}

	for _, tc := range tcc {
		layout := "2006-01-02T15:04:05.999999999Z"
		timeParsed, _ := time.Parse(layout, tc.input)
		timeParsed, _ = ParseTimeNanoToRFC3339(timeParsed)
		assert.Equal(t, tc.output, timeParsed.String())
	}
}
