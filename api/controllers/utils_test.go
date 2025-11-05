package controllers

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ParseRange(t *testing.T) {
	testCases := []struct {
		desc          string
		from          string
		to            string
		expectedFrom  time.Time
		expectedTo    time.Time
		expectedError error
	}{
		{
			desc:          "empty range",
			from:          "",
			to:            "",
			expectedFrom:  time.Time{},
			expectedTo:    time.Time{},
			expectedError: nil,
		},
		{
			desc:          "invalid from date format",
			from:          "2025-01-01 12:00:00",
			to:            "2025-01-02",
			expectedFrom:  time.Time{},
			expectedTo:    time.Time{},
			expectedError: fmt.Errorf("invalid from date format: the format must be YYYY-MM-DD"),
		},
		{
			desc:          "invalid to date format",
			from:          "2025-01-01",
			to:            "2025-01-02 12:00:00",
			expectedFrom:  time.Time{},
			expectedTo:    time.Time{},
			expectedError: fmt.Errorf("invalid to date format: the format must be YYYY-MM-DD"),
		},
		{
			desc:          "from is after to",
			from:          "2025-01-02",
			to:            "2025-01-01",
			expectedFrom:  time.Time{},
			expectedTo:    time.Time{},
			expectedError: fmt.Errorf("'To' date cannot be earlier than 'From' date"),
		},
		{
			desc:          "valid range",
			from:          "2025-01-01",
			to:            "2025-01-02",
			expectedFrom:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedTo:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedError: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			q := url.Values{}
			q.Add("from", tC.from)
			q.Add("to", tC.to)

			req := httptest.NewRequest("GET", "http://localhost:8080", nil)
			req.URL.RawQuery = q.Encode()

			from, to, err := parseDateRange(req)

			if tC.expectedError != nil {
				assert.Error(t, err, tC.expectedError.Error())
			}

			assert.Equal(t, tC.expectedFrom, from)
			assert.Equal(t, tC.expectedTo, to)

		})
	}
}
