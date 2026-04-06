package utils

import (
	"bufio"
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"social-network/internal/models"
	"social-network/internal/sessions"
)

//--------------------------------------------------------------------------------------|

// LoadEnv loads environment variables from a .env file at the specified path.
func LoadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}

		os.Setenv(key, val)
	}

	return scanner.Err()
}

//--------------------------------------------------------------------------------------|

// Getenv retrieves the value of an environment variable or returns a default value.
func Getenv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

//--------------------------------------------------------------------------------------|

// GetIntEnv retrieves the value of an environment variable as an integer or returns a default value.
func GetIntEnv(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

//--------------------------------------------------------------------------------------|

// GetUserID extracts the user ID from the request context or session cookie.
func GetUserID(ctx context.Context, r *http.Request, sm *sessions.SessionManager) int {
	if identity, ok := ctx.Value(models.UserKey).(*models.UserIdentity); ok {
		return identity.ID
	}

	if r == nil {
		return 0
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0
	}

	session, err := sm.GetSession(ctx, cookie.Value)
	if err != nil {
		return 0
	}
	return session.UserID
}

//--------------------------------------------------------------------------------------|

// IsLoggedIn checks if a user is currently logged in.
func IsLoggedIn(ctx context.Context, r *http.Request, sm *sessions.SessionManager) bool {
	return GetUserID(ctx, r, sm) > 0
}

//--------------------------------------------------------------------------------------|

// ParseInt is a helper to convert a string to an integer.
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
