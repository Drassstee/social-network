package user

import (
	"encoding/json"
	"net/http"

	"social-network/internal/models"
)

var errInternalServer = map[string]string{"error": "internal server error"}

type UserHandler struct {
	Service models.UserService
}

func NewUserHandler(us models.UserService) *UserHandler {
	return &UserHandler{Service: us}
}

type errJSON struct {
	Error string `json:"error"`
}

func notImplemented(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(errJSON{Error: "not implemented"})
}
