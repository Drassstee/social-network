package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"social-network/internal/models"
	"social-network/internal/models/user"
	"social-network/internal/utils"
	"social-network/internal/web"

	"github.com/google/uuid"
)



type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}



func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request, _ *models.UserIdentity) error {
	var u user.User

	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}

		u.FirstName = r.FormValue("first_name")
		u.LastName = r.FormValue("last_name")
		u.Email = r.FormValue("email")
		u.Password = r.FormValue("password")
		u.Nickname = r.FormValue("nickname")
		u.AboutMe = r.FormValue("about_me")
		dobStr := r.FormValue("dob")
		if dobStr != "" {
			dob, err := time.Parse(time.RFC3339, dobStr)
			if err == nil {
				u.DOB = &dob
			}
		}

		file, header, err := r.FormFile("avatar")
		if err == nil {
			defer file.Close()

			ext := filepath.Ext(header.Filename)
			dir := "uploads/avatars/registration"
			os.MkdirAll(dir, 0750)

			filename := uuid.NewString() + ext
			filePath := dir + "/" + filename
			dst, err := os.Create(filePath)
			if err == nil {
				defer dst.Close()
				io.Copy(dst, file)
				u.AvatarURL = filePath
			}
		}
	} else {
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			return web.StatusError{Code: http.StatusBadRequest, Err: err}
		}
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
