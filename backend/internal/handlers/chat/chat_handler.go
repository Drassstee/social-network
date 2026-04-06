package chat

import (
	"context"
	"errors"
	"io"
	"net/http"
	"social-network/internal/models"
	chatsvc "social-network/internal/service/chat"
	"social-network/internal/web"
)

//--------------------------------------------------------------------------------------|

// ChatHandler handles HTTP requests for chat messages and online users.
type ChatHandler struct {
	Repo     models.ChatRepo
	Hub      *chatsvc.Hub
	UserRepo models.UserRepo
	Uploader ImageUploader
}

//--------------------------------------------------------------------------------------|

// NewChatHandler creates a new instance of the chat handler.
func NewChatHandler(repo models.ChatRepo, hub *chatsvc.Hub, userRepo models.UserRepo, uploader ImageUploader) *ChatHandler {
	return &ChatHandler{
		Repo:     repo,
		Hub:      hub,
		UserRepo: userRepo,
		Uploader: uploader,
	}
}

//--------------------------------------------------------------------------------------|

// ImageUploader defines the interface for uploading chat images.
// This abstraction allows the ChatHandler to remain independent of the specific
// storage implementation (e.g., local disk, cloud storage).
type ImageUploader interface {
	UploadImage(ctx context.Context, userID int, filename string, content io.Reader) (string, error)
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
	otherUserID := web.QueryInt(r, "user_id", 0)
	if otherUserID == 0 {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("Invalid user ID")}
	}

	limit := web.QueryInt(r, "limit", 10)
	offset := web.QueryInt(r, "offset", 0)

	messages, err := h.Repo.GetMessages(r.Context(), identity.ID, otherUserID, limit, offset)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, messages)
	return nil
}

//--------------------------------------------------------------------------------------|

// GetOnlineUsers returns a list of all currently connected users.
func (h *ChatHandler) GetOnlineUsers(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	onlineIDs := h.Hub.GetOnlineUsers()

	users, err := h.UserRepo.GetByIDs(r.Context(), onlineIDs)
	if err != nil {
		return err
	}

	// Transform models.User to the simplified OnlineUser struct for API response
	type OnlineUser struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}

	onlineUsers := make([]OnlineUser, len(users))
	for i, u := range users {
		onlineUsers[i] = OnlineUser{
			ID:       u.ID,
			Username: u.Username,
		}
	}

	web.JSONResponse(w, http.StatusOK, onlineUsers)
	return nil
}
