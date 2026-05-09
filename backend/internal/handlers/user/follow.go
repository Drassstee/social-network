package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/web"
)

// --------------------------------------------------------------------|

func (h *UserHandler) Follow(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}

	var f models.Follow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	f.FollowerID = identity.ID

	msg, err := h.Service.Follow(f)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}
		return web.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"message": msg})
	return nil
}

// --------------------------------------------------------------------|

func (h *UserHandler) Unfollow(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}

	var f models.Follow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	f.FollowerID = identity.ID

	if err := h.Service.Unfollow(f); err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}
		return web.StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"message": "unfollowed successfully"})
	return nil
}
