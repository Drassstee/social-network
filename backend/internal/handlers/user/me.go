package user

import (
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/utils"
	"social-network/internal/web"
)

// GetMe returns the currently authenticated user's data.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}

	id := int64(identity.ID)
	data, err := h.Users.GetMe(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return web.StatusError{Code: http.StatusNotFound, Err: err}
		}
		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, data)
	return nil
}
