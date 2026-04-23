package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/models/user"
	"social-network/internal/utils"
	"social-network/internal/web"
)

// --------------------------------------------------------------------|

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	userID := int64(identity.ID)

	targetID, err := utils.GetIDByURL(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	u, err := h.Users.GetProfile(userID, targetID)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		} else if errors.Is(err, models.ErrNotFound) {
			return web.StatusError{Code: http.StatusNotFound, Err: err}
		} else if errors.Is(err, models.ErrUserPrivate) {
			return web.StatusError{Code: http.StatusForbidden, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, u)
	return nil
}

// --------------------------------------------------------------------|

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	id := int64(identity.ID)

	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	u.ID = id

	err = h.Users.UpdateProfile(&u)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		} else if errors.Is(err, models.ErrConflict) {
			return web.StatusError{Code: http.StatusConflict, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, u)
	return nil
}
