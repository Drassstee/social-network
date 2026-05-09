package chat

import (
	"errors"
	"net/http"
	"social-network/internal/models"
	chatsvc "social-network/internal/service/chat"
	"social-network/internal/web"
)

//--------------------------------------------------------------------------------------|

// ChatHandler handles HTTP requests for chat messages and online users.
type ChatHandler struct {
	Service  models.ChatService
	Hub      *chatsvc.Hub
	Uploader models.ImageUploader
}

//--------------------------------------------------------------------------------------|

// NewChatHandler creates a new instance of the chat handler.
func NewChatHandler(service models.ChatService, hub *chatsvc.Hub, uploader models.ImageUploader) *ChatHandler {
	return &ChatHandler{
		Service:  service,
		Hub:      hub,
		Uploader: uploader,
	}
}

//--------------------------------------------------------------------------------------|

// UploadImage handles the HTTP request to upload an image for a chat message.
func (h *ChatHandler) UploadImage(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if r.Method != http.MethodPost {
		return web.StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("Method not allowed")}
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("Missing image file")}
	}
	defer file.Close()

	if header.Size > web.MaxImageSize {
		return web.StatusError{Code: http.StatusBadRequest, Err: models.ErrImageTooBig}
	}

	url, err := h.Uploader.UploadImage(r.Context(), identity.ID, header.Filename, file)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"url": url})
	return nil
}

//--------------------------------------------------------------------------------------|

// GetMessages returns a list of private messages between two users.
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	otherUserID := web.QueryInt64(r, "user_id", 0)
	if otherUserID == 0 {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("Invalid user ID")}
	}

	limit := web.QueryInt(r, "limit", 10)
	offset := web.QueryInt(r, "offset", 0)

	messages, err := h.Service.GetChatHistory(r.Context(), identity.ID, otherUserID, limit, offset)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, messages)
	return nil
}

//--------------------------------------------------------------------------------------|

// GetChatableUsers returns a list of all users the current user can chat with, along with their online status.
func (h *ChatHandler) GetChatableUsers(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	users, err := h.Service.GetChatableUsers(r.Context(), identity.ID)
	if err != nil {
		return err
	}

	onlineIDs := h.Hub.GetOnlineUsers()
	onlineMap := make(map[int64]bool)
	for _, id := range onlineIDs {
		onlineMap[id] = true
	}

	// Enrich with online status
	type enrichedUser struct {
		models.OnlineUser
		IsOnline bool `json:"is_online"`
	}

	enriched := make([]enrichedUser, len(users))
	for i, u := range users {
		enriched[i] = enrichedUser{
			OnlineUser: u,
			IsOnline:   onlineMap[u.ID],
		}
	}

	web.JSONResponse(w, http.StatusOK, enriched)
	return nil
}

//--------------------------------------------------------------------------------------|

// GetUnreadCounts returns the unread message counts for the authenticated user.
func (h *ChatHandler) GetUnreadCounts(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	counts, err := h.Service.GetUnreadCounts(r.Context(), identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, counts)
	return nil
}

//--------------------------------------------------------------------------------------|

// MarkAsRead marks messages from a specific sender as read for the current user.
func (h *ChatHandler) MarkAsRead(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	senderID := web.QueryInt64(r, "sender_id", 0)
	if senderID == 0 {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("Invalid sender ID")}
	}

	if err := h.Service.MarkAsRead(r.Context(), senderID, identity.ID); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	return nil
}

//--------------------------------------------------------------------------------------|

// Connect upgrades the HTTP connection to a WebSocket connection.
func (h *ChatHandler) Connect(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	chatsvc.ServeWs(h.Hub, w, r, identity.ID)
	return nil
}

