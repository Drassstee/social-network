package user

import (
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/web"
)

// GetMe returns the currently authenticated user's data.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}

	data, err := h.Service.GetMe(identity.ID)
	if err != nil {
		return web.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	web.JSONResponse(w, http.StatusOK, data)
	return nil
}
