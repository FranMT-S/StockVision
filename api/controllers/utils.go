package controllers

import (
	"api/models"
	"api/models/filters"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// parseFilters extracts pagination and ordering parameters from query string
func parseFilters(r *http.Request) filters.Filters {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	size, _ := strconv.Atoi(query.Get("size"))
	sortStr := strings.ToLower(query.Get("sort"))
	queryStr := query.Get("q")

	var sort filters.Sort
	switch sortStr {
	case "asc":
		sort = filters.ASC
	case "desc":
		sort = filters.DESC
	default:
		sort = filters.ASC
	}

	return filters.Filters{
		Page:     page,
		PageSize: size,
		Sort:     sort,
		Query:    queryStr,
	}
}

// parseDateRange extracts date range parameters from query string
// validate the format is correct
// validate that to is not before from
func parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	var err error

	var fromTime time.Time
	if from != "" {
		fromTime, err = time.Parse("2006-01-02", from)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date format: the format must be YYYY-MM-DD")
		}
	}

	var toTime time.Time
	if to != "" {
		toTime, err = time.Parse("2006-01-02", to)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to date format: the format must be YYYY-MM-DD")
		}

		if fromTime.After(toTime) {
			return time.Time{}, time.Time{}, fmt.Errorf("'To' date cannot be earlier than 'From' date")
		}
	}

	return fromTime, toTime, nil
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.NewResponseError(message))
}
