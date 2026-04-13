package user

import (
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/models/avatar"
	"social-network/internal/utils"
)

func (h *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	userID, err := utils.GetIDByURL(r, "id")
	if err != nil {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}

	avatarURL, err := h.Users.GetAvatar(userID)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusBadRequest, msg)
			return
		} else if errors.Is(err, models.ErrNotFound) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusNotFound, msg)
			return
		}
		utils.RespondJSON(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	http.ServeFile(w, r, avatarURL)
}

// --------------------------------------------------------------------|

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		msg := map[string]string{"error": "file too large"}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
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
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusBadRequest, msg)
			return
		}

		utils.RespondJSON(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	utils.RespondJSON(w, http.StatusOK, nil)
}
