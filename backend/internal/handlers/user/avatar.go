package user

import (
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/models/avatar"
	"social-network/internal/utils"
	"social-network/internal/web"
)

func (h *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}

	userID, err := utils.GetIDByURL(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	avatarURL, err := h.Users.GetAvatar(userID)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		} else if errors.Is(err, models.ErrNotFound) {
			return web.StatusError{Code: http.StatusNotFound, Err: err}
		}
		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	http.ServeFile(w, r, avatarURL)
	return nil
}

// --------------------------------------------------------------------|

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	userID := int64(identity.ID)

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("file too large")}
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("invalid file")}
	}
	defer file.Close()

	var a = avatar.Avatar{
		File:   file,
		Header: header,
		UserID: userID,
	}

	err = h.Users.UploadAvatar(&a)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusOK, nil)
	return nil
}
