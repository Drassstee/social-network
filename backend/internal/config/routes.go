package config

import (
	"net/http"

	"social-network/internal/handlers"
	"social-network/internal/middleware"
)

func SetupRoutes(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	apimux := http.NewServeMux()

	apimux.HandleFunc("POST /register", h.Register)
	apimux.HandleFunc("POST /login", h.Login)

	apimux.HandleFunc("POST /logout", middleware.AuthMiddleware(h.Users, h.Logout))
	apimux.HandleFunc("DELETE /delete", middleware.AuthMiddleware(h.Users, h.DeleteUser))
	apimux.HandleFunc("GET /users/{id}", middleware.AuthMiddleware(h.Users, h.GetProfile))
	apimux.HandleFunc("PUT /users", middleware.AuthMiddleware(h.Users, h.UpdateProfile))

	apimux.HandleFunc("POST /follow", middleware.AuthMiddleware(h.Users, h.Follow))
	apimux.HandleFunc("POST /unfollow", middleware.AuthMiddleware(h.Users, h.Unfollow))

	apimux.HandleFunc("GET /notifications", middleware.AuthMiddleware(h.Users, h.GetNotification))
	apimux.HandleFunc("POST /notifications", middleware.AuthMiddleware(h.Users, h.RespondToFollowRequest))

	apimux.HandleFunc("GET /avatar/{id}", middleware.AuthMiddleware(h.Users, h.GetAvatar))
	apimux.HandleFunc("POST /avatar", middleware.AuthMiddleware(h.Users, h.UploadAvatar))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apimux))

	return mux
}
