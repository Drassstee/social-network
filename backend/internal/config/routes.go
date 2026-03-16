package config

import (
	"net/http"
	"social-network/internal/handlers"
)

func SetupRoutes(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	apimux := http.NewServeMux()
	// этих роутов все еще нет
	apimux.HandleFunc("POST /register", h.User.Register)
	apimux.HandleFunc("POST /login", h.User.Login)
	apimux.HandleFunc("POST /logout", h.User.Logout)
	apimux.HandleFunc("GET /users/{id}", h.User.GetProfile)
	apimux.HandleFunc("GET /posts", h.Post.GetPosts)
	mux.Handle("/api/v1", http.StripPrefix("api/v1", apimux))

	return mux
}
