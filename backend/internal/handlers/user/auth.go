package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"social-network/internal/models/user"
	"social-network/internal/utils"
)

// --------------------------------------------------------------------|

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// --------------------------------------------------------------------|

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		msg := map[string]string{"error": err.Error()} // invalid json data
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}

	data, err := h.Users.Register(&u)
	if errors.Is(err, user.ErrExists) {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusConflict, msg)
		return
	} else if err != nil {
		msg := map[string]string{"error": "internal server error"}
		utils.RespondJSON(w, http.StatusInternalServerError, msg)
		return
	}

	utils.SetCookie(w, data.UUID, data.ExpiresAt)

	fmt.Println("User created")
	utils.RespondJSON(w, http.StatusCreated, data)
}

// --------------------------------------------------------------------|

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var d LoginData
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		msg := map[string]string{"error": err.Error()} // invalid json data
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	}

	data, err := h.Users.Login(d.Email, d.Password)
	if errors.Is(err, user.ErrInvalidData) {
		msg := map[string]string{"error": err.Error()}
		utils.RespondJSON(w, http.StatusBadRequest, msg)
		return
	} else if err != nil {
		msg := map[string]string{"error": "internal server error"}
		utils.RespondJSON(w, http.StatusInternalServerError, msg)
		return
	}

	utils.SetCookie(w, data.UUID, data.ExpiresAt)

	fmt.Println("The user has logged in")
	utils.RespondJSON(w, http.StatusOK, data)
}

// --------------------------------------------------------------------|

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetUserIDByContext(r)
	if !ok {
		msg := map[string]string{"error": "unauthorized"}
		utils.RespondJSON(w, http.StatusUnauthorized, msg)
		return
	}

	if err := h.Users.Logout(id); err != nil {
		msg := map[string]string{"error": "internal server error"}
		utils.RespondJSON(w, http.StatusInternalServerError, msg)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "the user has logged out"})
}

// --------------------------------------------------------------------|
