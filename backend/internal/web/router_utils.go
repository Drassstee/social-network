package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

//--------------------------------------------------------------------------------------|

// ParseJSON decodes the request body into the provided data structure.
func ParseJSON(r *http.Request, data any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return StatusError{Code: http.StatusBadRequest, Err: errors.New("Content-Type must be application/json")}
	}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return StatusError{Code: http.StatusBadRequest, Err: errors.New("Invalid JSON")}
	}
	return nil
}

//--------------------------------------------------------------------------------------|

// QueryInt extracts an integer from the URL query parameters with a default value.
func QueryInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

//--------------------------------------------------------------------------------------|

// FormInt extracts an integer from the request form (query or body) with a default value.
func FormInt(r *http.Request, key string, defaultVal int) int {
	val := r.FormValue(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

//--------------------------------------------------------------------------------------|

// RouteInt extracts an integer from the path using the provided segment index.
func RouteInt(r *http.Request, segmentIndex int) (int, error) {
	return ExtractIDFromPath(r.URL.Path, segmentIndex)
}
