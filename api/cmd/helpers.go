package cmd

import "time"

// ParseTimeNanoToRFC3339 parses a time string in RFC3339Nano format
// 2025-09-15T00:30:05.082205952Z
func ParseTimeNanoToRFC3339(t time.Time) (time.Time, error) {
	return time.Parse(time.RFC3339, t.Format(time.RFC3339))
}
