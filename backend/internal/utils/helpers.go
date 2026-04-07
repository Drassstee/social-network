package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"social-network/internal/models/user"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "invalid json", http.StatusInternalServerError)
		return
	}
}

// --------------------------------------------------------------------|

func SetCookie(w http.ResponseWriter, uuid string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    uuid,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

// --------------------------------------------------------------------|

func DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// --------------------------------------------------------------------|

func GetUserIDByContext(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(user.Key).(int64)
	return userID, ok
}

// --------------------------------------------------------------------|

func GetIDByURL(r *http.Request, target string) (int64, error) {
	idStr := r.PathValue(target)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
