package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"social-network/internal/models"
	"social-network/internal/models/user"
	"social-network/internal/utils"
	"social-network/internal/web"
)



type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}



func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request, _ *models.UserIdentity) error {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	data, err := h.Users.Register(&u)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		} else if errors.Is(err, models.ErrConflict) {
			return web.StatusError{Code: http.StatusConflict, Err: err}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.SetCookie(w, data.UUID, *data.ExpiresAt)

	fmt.Println("User created")
	utils.RespondJSON(w, http.StatusCreated, data)
	return nil
}



func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request, _ *models.UserIdentity) error {
	var d LoginData
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	data, err := h.Users.Login(d.Email, d.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("invalid email or password")}
		}

		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.SetCookie(w, data.UUID, *data.ExpiresAt)

	fmt.Println("The user has logged in")
	utils.RespondJSON(w, http.StatusOK, data)
	return nil
}



func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	id := int64(identity.ID)

	err := h.Users.Logout(id)
	if err != nil {
		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.DeleteCookie(w)

	utils.RespondJSON(w, http.StatusNoContent, nil)
	return nil
}



func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	id := int64(identity.ID)

	err := h.Users.DeleteUser(id)
	if err != nil {
		if errors.Is(err, models.ErrInvalidData) {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}
		return web.StatusError{Code: http.StatusInternalServerError, Err: errors.New("internal server error")}
	}

	utils.RespondJSON(w, http.StatusNoContent, nil)
	return nil
}
