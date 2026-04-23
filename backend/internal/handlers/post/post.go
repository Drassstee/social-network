package post

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"social-network/internal/models"
	servicepost "social-network/internal/service/post"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, userID int, filename string, content io.Reader) (string, error)
}

type PostHandler struct {
	serv     *servicepost.PostService
	uploader ImageUploader
}

func NewPostHandler(serv *servicepost.PostService, uploader ImageUploader) *PostHandler {
	return &PostHandler{serv: serv, uploader: uploader}
}

type postsListResponse struct {
	Posts   []models.Post `json:"posts"`
	HasMore bool          `json:"has_more"`
}

type createPostRequest struct {
	AuthorID string `json:"author_id"`
	Content  string `json:"content"`
}

type createPostResponse struct {
	Post *models.Post `json:"post"`
}

type errResponse struct {
	Error string `json:"error"`
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	groupID, _ := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	
	requesterID := int64(0)
	if identity != nil {
		requesterID = int64(identity.ID)
	}

	posts, hasMore, err := h.serv.ListPosts(requesterID, groupID, limit, offset)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, postsListResponse{Posts: posts, HasMore: hasMore})
	return nil
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return &models.ValidationError{Field: "auth", Message: "unauthorized"}
	}

	// Support multi-part for image uploads
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return h.handleMultipartCreatePost(w, r, identity)
	}

	var body struct {
		Content      string  `json:"content"`
		Privacy      string  `json:"privacy"`
		GroupID      int64   `json:"group_id"`
		AllowedUsers []int64 `json:"allowed_users"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return &models.ValidationError{Field: "json", Message: "invalid json"}
	}

	p, err := h.serv.CreatePost(int64(identity.ID), body.Content, "", body.Privacy, body.GroupID, body.AllowedUsers)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusCreated, createPostResponse{Post: p})
	return nil
}

func (h *PostHandler) handleMultipartCreatePost(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	err := r.ParseMultipartForm(20 << 20) // 20MB
	if err != nil {
		return err
	}

	content := r.FormValue("content")
	privacy := r.FormValue("privacy")
	allowedUsersStr := r.FormValue("allowed_users")
	var allowedUsers []int64
	if allowedUsersStr != "" {
		json.Unmarshal([]byte(allowedUsersStr), &allowedUsers)
	}

	groupID, _ := strconv.ParseInt(r.FormValue("group_id"), 10, 64)

	var imageURL string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageURL, err = h.uploader.UploadImage(r.Context(), identity.ID, header.Filename, file)
		if err != nil {
			return err
		}
	}

	p, err := h.serv.CreatePost(int64(identity.ID), content, imageURL, privacy, groupID, allowedUsers)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusCreated, createPostResponse{Post: p})
	return nil
}

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return &models.ValidationError{Field: "auth", Message: "unauthorized"}
	}

	var content, imageURL string
	var postID int64

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			return err
		}
		content = r.FormValue("content")
		postID, _ = strconv.ParseInt(r.FormValue("post_id"), 10, 64)
		
		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imageURL, err = h.uploader.UploadImage(r.Context(), identity.ID, header.Filename, file)
			if err != nil {
				return err
			}
		}
	} else {
		var body struct {
			PostID  int64  `json:"post_id"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return err
		}
		content = body.Content
		postID = body.PostID
	}

	c, err := h.serv.CreateComment(int64(identity.ID), postID, content, imageURL)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusCreated, c)
	return nil
}


func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
