package utils

import (
	"os"
	"strconv"
	"strings"
)

// Getenv retrieves the value of an environment variable or returns a default value.
func Getenv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// ParseInt is a helper to convert a string to an int64.
func ParseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// FormatAvatarURL converts a disk path to a web URL.
func FormatAvatarURL(path string) string {
	if path == "" {
		return ""
	}
	// If it's already a full URL or starts with /, return as is
	if path[0] == '/' || strings.HasPrefix(path, "http") {
		return path
	}
	// Prefix with api path
	return "/api/v1/" + path
}
