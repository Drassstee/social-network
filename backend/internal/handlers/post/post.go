package post

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errResponse{Error: "method not allowed"})
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	posts, hasMore, err := h.serv.ListPosts(limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResponse{Error: "failed to load posts"})
		return
	}
	writeJSON(w, http.StatusOK, postsListResponse{Posts: posts, HasMore: hasMore})
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errResponse{Error: "method not allowed"})
		return
	}
	var body createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResponse{Error: "invalid json"})
		return
	}
	p, err := h.serv.CreatePost(body.AuthorID, body.Content)
	if err != nil {
		msg := err.Error()
		if strings.Contains(strings.ToUpper(msg), "FOREIGN KEY") {
			writeJSON(w, http.StatusBadRequest, errResponse{Error: "author_id does not exist"})
			return
		}
		writeJSON(w, http.StatusBadRequest, errResponse{Error: msg})
		return
	}
	writeJSON(w, http.StatusCreated, createPostResponse{Post: p})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
