package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/models/user"
	"social-network/internal/utils"
)

// --------------------------------------------------------------------|

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	targetID, err := utils.GetIDByURL(r, "id")
	if err != nil {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}

	u, err := h.Users.GetProfile(userID, targetID)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusBadRequest, msg)
			return
		} else if errors.Is(err, models.ErrNotFound) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusNotFound, msg)
			return
		} else if errors.Is(err, models.ErrUserPrivate) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusForbidden, msg)
			return
		}

		utils.RespondJSON(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	utils.RespondJSON(w, http.StatusOK, u)
}

// --------------------------------------------------------------------|

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}
	u.ID = id

	err = h.Users.UpdateProfile(&u)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusBadRequest, msg)
			return
		} else if errors.Is(err, models.ErrConflict) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusConflict, msg)
			return
		}

		utils.RespondJSON(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	utils.RespondJSON(w, http.StatusOK, u)
}
