package post

import (
	"encoding/json"
	"net/http"
	"strconv"

	"social-network/internal/models"
	servicepost "social-network/internal/service/post"
)

type PostHandler struct {
	serv *servicepost.PostService
}

func NewPostHandler(serv *servicepost.PostService) *PostHandler {
	return &PostHandler{serv: serv}
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
	posts, hasMore, err := h.serv.ListPosts(limit, offset)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, postsListResponse{Posts: posts, HasMore: hasMore})
	return nil
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	var body createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return &models.ValidationError{Field: "json", Message: "invalid json"}
	}
	p, err := h.serv.CreatePost(body.AuthorID, body.Content)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusCreated, createPostResponse{Post: p})
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
