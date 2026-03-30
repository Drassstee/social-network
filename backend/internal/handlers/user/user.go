package user

import (
	"encoding/json"
	"net/http"

	"social-network/internal/models"
)

type UserHandler struct {
	serv models.UserService
}

func NewUserHandler(serv models.UserService) *UserHandler {
	return &UserHandler{serv: serv}
}

type errJSON struct {
	Error string `json:"error"`
}

func notImplemented(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(errJSON{Error: "not implemented"})
}

// Stubs until auth/profile is implemented by teammate.

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}
