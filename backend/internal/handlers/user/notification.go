package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/web"
)

// --------------------------------------------------------------------|

func (h *UserHandler) RespondToFollowRequest(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}

	var f models.Follow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	f.FollowingID = identity.ID

	err := h.Service.RespondToFollowRequest(f)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}
		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"message": "request handled successfully"})
	return nil
}

// --------------------------------------------------------------------|
