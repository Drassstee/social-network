package utils

import (
	"os"
	"strconv"
	"strings"
)

//--------------------------------------------------------------------------------------|

func Getenv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

//--------------------------------------------------------------------------------------|

func ParseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

//--------------------------------------------------------------------------------------|

func FormatAvatarURL(path string) string {
	if path == "" {
		return ""
	}
	if path[0] == '/' || strings.HasPrefix(path, "http") {
		return path
	}
	return "/api/v1/" + path
}
