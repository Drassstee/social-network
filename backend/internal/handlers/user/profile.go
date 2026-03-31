package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"social-network/internal/models/user"
	"social-network/internal/utils"
)

// --------------------------------------------------------------------|

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	id, err := utils.GetIDByURL(r, "id") //<==

	u, err := h.Users.GetProfile(id)
	if errors.Is(err, user.ErrNotFound) {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusNotFound, msg)
		return
	} else if err != nil {
		msg := map[string]string{"error": "internal server error"}
		utils.RespondJSON(w, http.StatusInternalServerError, msg)
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
	if strings.Contains(err.Error(), "sql") {
		msg := map[string]string{"error": "internal server error"}
		utils.RespondJSON(w, http.StatusInternalServerError, msg)
		return
	} else if err != nil {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}

	utils.RespondJSON(w, http.StatusOK, u)
}
