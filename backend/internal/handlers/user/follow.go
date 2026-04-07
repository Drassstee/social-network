package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/models/follow"
	"social-network/internal/utils"
)

// --------------------------------------------------------------------|

func (h *UserHandler) Follow(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorization"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	var data follow.Follow
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}
	data.FollowerID = userID

	status, err := h.Users.Follow(data)
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

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": status})
}

// --------------------------------------------------------------------|

func (h *UserHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorization"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	var data follow.Follow
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}
	data.FollowerID = userID

	err = h.Users.Unfollow(data)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusBadRequest, msg)
			return
		}

		utils.RespondJSON(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "unfollow"})
}

// --------------------------------------------------------------------|
