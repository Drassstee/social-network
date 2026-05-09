package middleware

import (
	"context"
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/web"
)

//--------------------------------------------------------------------------------------|

func AuthMiddleware(serv models.UserService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) || cookie.Value == "" {
			msg := map[string]string{"error": "no cookie"}
			web.JSONResponse(w, http.StatusUnauthorized, msg)
			return
		}

		id, err := serv.GetUserID(cookie.Value)
		if errors.Is(err, models.ErrNotFound) {
			msg := map[string]string{"error": "unauthorized"}
			web.JSONResponse(w, http.StatusUnauthorized, msg)
			return
		} else if err != nil {
			msg := map[string]string{"error": "internal server error"}
			web.JSONResponse(w, http.StatusInternalServerError, msg)
			return
		}

		identity := &models.UserIdentity{
			ID: id,
		}

		ctx := context.WithValue(r.Context(), models.UserKey, identity)
		next(w, r.WithContext(ctx))
	}
}
