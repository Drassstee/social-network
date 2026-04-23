package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/models/follow"
	"social-network/internal/utils"
	"social-network/internal/web"
)

// --------------------------------------------------------------------|

func (h *UserHandler) Follow(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	userID := int64(identity.ID)

	var data follow.Follow
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	data.FollowerID = userID

	status, err := h.Users.Follow(data)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		} else if errors.Is(err, models.ErrConflict) {
			return web.StatusError{Code: http.StatusConflict, Err: err}
		} else if errors.Is(err, models.ErrNotFound) {
			return web.StatusError{Code: http.StatusNotFound, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": status})
	return nil
}

// --------------------------------------------------------------------|

func (h *UserHandler) Unfollow(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	userID := int64(identity.ID)

	var data follow.Follow
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	data.FollowerID = userID

	err = h.Users.Unfollow(data)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "unfollow"})
	return nil
}

// --------------------------------------------------------------------|
