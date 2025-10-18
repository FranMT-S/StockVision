package filters

import "strings"

type Sort string

const (
	ASC  Sort = "asc"
	DESC Sort = "desc"
)

// String returns the string representation of the sort, defaulting to "asc" if the sort is invalid
func (o Sort) String() string {
	if !o.IsValid() {
		return "asc"
	}

	return strings.ToLower(string(o))
}

func (o Sort) IsValid() bool {
	return o == ASC || o == DESC
}
