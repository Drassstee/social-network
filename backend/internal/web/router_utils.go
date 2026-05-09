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
	return int(QueryInt64(r, key, int64(defaultVal)))
}

//--------------------------------------------------------------------------------------|

// QueryInt64 extracts an int64 from the URL query parameters with a default value.
func QueryInt64(r *http.Request, key string, defaultVal int64) int64 {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return i
}

//--------------------------------------------------------------------------------------|

// PathInt64 extracts an int64 from the URL path using Go 1.22 path values.
func PathInt64(r *http.Request, key string) (int64, error) {
	val := r.PathValue(key)
	if val == "" {
		return 0, errors.New("missing path value")
	}
	return strconv.ParseInt(val, 10, 64)
}
