package config

import (
	"net/http"
	"social-network/internal/handlers"
	"social-network/internal/web"
)

func SetupRoutes(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	apimux := http.NewServeMux()
	// Auth
	apimux.HandleFunc("POST /register", h.User.Register)
	apimux.HandleFunc("POST /login", h.User.Login)
	apimux.HandleFunc("POST /logout", h.User.Logout)

	// Users
	apimux.HandleFunc("GET /users/{id}", h.User.GetProfile)

	// Posts
	apimux.HandleFunc("GET /posts", web.NewAppHandler(h.Post.GetPosts))
	apimux.HandleFunc("POST /posts", web.NewAppHandler(h.Post.CreatePost))

	// Groups
	apimux.HandleFunc("POST /groups", web.NewAppHandler(h.Group.CreateGroup))
	apimux.HandleFunc("GET /groups", web.NewAppHandler(h.Group.ListGroups))
	apimux.HandleFunc("GET /groups/{id}", web.NewAppHandler(h.Group.GetGroup))
	apimux.HandleFunc("GET /groups/{id}/members", web.NewAppHandler(h.Group.GetMembers))
	apimux.HandleFunc("POST /groups/{id}/leave", web.NewAppHandler(h.Group.LeaveGroup))
	apimux.HandleFunc("POST /groups/{id}/invite", web.NewAppHandler(h.Group.InviteUser))
	apimux.HandleFunc("GET /groups/invitations", web.NewAppHandler(h.Group.GetPendingInvitations))
	apimux.HandleFunc("POST /groups/invitations/{id}/respond", web.NewAppHandler(h.Group.RespondToInvitation))
	apimux.HandleFunc("POST /groups/{id}/request", web.NewAppHandler(h.Group.RequestJoin))
	apimux.HandleFunc("GET /groups/{id}/requests", web.NewAppHandler(h.Group.GetPendingJoinRequests))
	apimux.HandleFunc("POST /groups/requests/{id}/respond", web.NewAppHandler(h.Group.RespondToJoinRequest))
	apimux.HandleFunc("POST /groups/{id}/events", web.NewAppHandler(h.Group.CreateEvent))
	apimux.HandleFunc("GET /groups/{id}/events", web.NewAppHandler(h.Group.GetGroupEvents))
	apimux.HandleFunc("POST /groups/events/{id}/respond", web.NewAppHandler(h.Group.RespondToEvent))
	apimux.HandleFunc("GET /groups/{id}/messages", web.NewAppHandler(h.Group.GetGroupMessages))

	// Chat
	apimux.HandleFunc("GET /chat/messages", web.NewAppHandler(h.Chat.GetMessages))
	apimux.HandleFunc("POST /chat/upload", web.NewAppHandler(h.Chat.UploadImage))
	apimux.HandleFunc("GET /chat/online", web.NewAppHandler(h.Chat.GetOnlineUsers))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apimux))

	return mux
}
