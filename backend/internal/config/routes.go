package config

import (
	"net/http"

	"social-network/internal/handlers"
	"social-network/internal/middleware"
	"social-network/internal/web"
)

func SetupRoutes(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	apimux := http.NewServeMux()

	authMW := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(h.User.Users, next)
	}

	// Auth
	apimux.HandleFunc("POST /register", web.NewAppHandler(h.User.Register))
	apimux.HandleFunc("POST /login", web.NewAppHandler(h.User.Login))
	apimux.HandleFunc("POST /logout", authMW(web.NewAppHandler(h.User.Logout)))

	// Users
	apimux.HandleFunc("GET /users/{id}", authMW(web.NewAppHandler(h.User.GetProfile)))
	apimux.HandleFunc("DELETE /delete", authMW(web.NewAppHandler(h.User.DeleteUser)))
	apimux.HandleFunc("PUT /users", authMW(web.NewAppHandler(h.User.UpdateProfile)))
	apimux.HandleFunc("GET /notifications", authMW(web.NewAppHandler(h.User.GetNotification)))
	apimux.HandleFunc("POST /follow", authMW(web.NewAppHandler(h.User.Follow)))
	apimux.HandleFunc("POST /unfollow", authMW(web.NewAppHandler(h.User.Unfollow)))
	apimux.HandleFunc("POST /avatar", authMW(web.NewAppHandler(h.User.UploadAvatar)))
	apimux.HandleFunc("GET /avatar/{id}", authMW(web.NewAppHandler(h.User.GetAvatar)))

	// Posts
	apimux.HandleFunc("GET /posts", authMW(web.NewAppHandler(h.Post.GetPosts)))
	apimux.HandleFunc("POST /posts", authMW(web.NewAppHandler(h.Post.CreatePost)))

	// Groups
	apimux.HandleFunc("POST /groups", authMW(web.NewAppHandler(h.Group.CreateGroup)))
	apimux.HandleFunc("GET /groups", authMW(web.NewAppHandler(h.Group.ListGroups)))
	apimux.HandleFunc("GET /groups/{id}", authMW(web.NewAppHandler(h.Group.GetGroup)))
	apimux.HandleFunc("GET /groups/{id}/members", authMW(web.NewAppHandler(h.Group.GetMembers)))
	apimux.HandleFunc("POST /groups/{id}/leave", authMW(web.NewAppHandler(h.Group.LeaveGroup)))
	apimux.HandleFunc("POST /groups/{id}/invite", authMW(web.NewAppHandler(h.Group.InviteUser)))
	apimux.HandleFunc("GET /groups/invitations", authMW(web.NewAppHandler(h.Group.GetPendingInvitations)))
	apimux.HandleFunc("POST /groups/invitations/{id}/respond", authMW(web.NewAppHandler(h.Group.RespondToInvitation)))
	apimux.HandleFunc("POST /groups/{id}/request", authMW(web.NewAppHandler(h.Group.RequestJoin)))
	apimux.HandleFunc("GET /groups/{id}/requests", authMW(web.NewAppHandler(h.Group.GetPendingJoinRequests)))
	apimux.HandleFunc("POST /groups/requests/{id}/respond", authMW(web.NewAppHandler(h.Group.RespondToJoinRequest)))
	apimux.HandleFunc("POST /groups/{id}/events", authMW(web.NewAppHandler(h.Group.CreateEvent)))
	apimux.HandleFunc("GET /groups/{id}/events", authMW(web.NewAppHandler(h.Group.GetGroupEvents)))
	apimux.HandleFunc("POST /groups/events/{id}/respond", authMW(web.NewAppHandler(h.Group.RespondToEvent)))
	apimux.HandleFunc("GET /groups/{id}/messages", authMW(web.NewAppHandler(h.Group.GetGroupMessages)))

	// Chat
	apimux.HandleFunc("GET /chat/messages", authMW(web.NewAppHandler(h.Chat.GetMessages)))
	apimux.HandleFunc("POST /chat/upload", authMW(web.NewAppHandler(h.Chat.UploadImage)))
	apimux.HandleFunc("GET /chat/online", authMW(web.NewAppHandler(h.Chat.GetOnlineUsers)))
	apimux.HandleFunc("GET /ws", authMW(web.NewAppHandler(h.Chat.Connect)))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apimux))

	return mux
}
