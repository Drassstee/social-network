package middleware

import (
	"context"
	"errors"
	"net/http"

	Serv "social-network/internal/handlers/user"
	"social-network/internal/models"
	"social-network/internal/models/user"
	"social-network/internal/utils"
)

func AuthMiddleware(serv Serv.UserServ, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) || cookie.Value == "" {
			msg := map[string]string{"error": "no cookie"}
			utils.RespondJSON(w, http.StatusUnauthorized, msg)
			return
		}

		id, err := serv.GetUserID(cookie.Value)
		if errors.Is(err, models.ErrNotFound) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusUnauthorized, msg)
			return
		} else if err != nil {
			msg := map[string]string{"error": "internal server error"}
			utils.RespondJSON(w, http.StatusInternalServerError, msg)
			return
		}

		ctx := context.WithValue(r.Context(), user.Key, id)
		next(w, r.WithContext(ctx))
	}
}
