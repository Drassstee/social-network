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

func (h *UserHandler) GetNotification(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	id := int64(identity.ID)

	users, err := h.Users.GetNotification(id)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, users)
	return nil
}

// --------------------------------------------------------------------|

func (h *UserHandler) RespondToFollowRequest(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	id := int64(identity.ID)

	var f follow.Follow
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	f.FollowingID = id

	err = h.Users.RespondToFollowRequest(f)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		} else if errors.Is(err, models.ErrNotFound) {
			return web.StatusError{Code: http.StatusNotFound, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, nil)
	return nil
}

// --------------------------------------------------------------------|
