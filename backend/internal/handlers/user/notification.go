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

func (h *UserHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	users, err := h.Users.GetNotification(id)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			msg := map[string]string{"error": err.Error()}
			utils.RespondJSON(w, http.StatusBadRequest, msg)
			return
		}

		utils.RespondJSON(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

// --------------------------------------------------------------------|

func (h *UserHandler) RespondToFollowRequest(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	var f follow.Follow
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}

	f.FollowingID = id

	err = h.Users.RespondToFollowRequest(f)
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

	utils.RespondJSON(w, http.StatusOK, nil)
}

// --------------------------------------------------------------------|
