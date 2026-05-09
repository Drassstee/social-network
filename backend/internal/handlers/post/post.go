package post

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"social-network/internal/models"
	"social-network/internal/web"
)

//--------------------------------------------------------------------------------------|

type PostHandler struct {
	serv     models.PostService
	uploader models.ImageUploader
}

//--------------------------------------------------------------------------------------|

func NewPostHandler(serv models.PostService, uploader models.ImageUploader) *PostHandler {
	return &PostHandler{serv: serv, uploader: uploader}
}

//--------------------------------------------------------------------------------------|

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

//--------------------------------------------------------------------------------------|

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	limit := web.QueryInt(r, "limit", 10)
	offset := web.QueryInt(r, "offset", 0)
	groupID, _ := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	
	requesterID := int64(0)
	if identity != nil {
		requesterID = identity.ID
	}

	posts, hasMore, err := h.serv.ListPosts(requesterID, groupID, limit, offset)
	if err != nil {
		return err
	}
	web.JSONResponse(w, http.StatusOK, postsListResponse{Posts: posts, HasMore: hasMore})
	return nil
}

//--------------------------------------------------------------------------------------|

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

	p, err := h.serv.CreatePost(identity.ID, body.Content, "", body.Privacy, body.GroupID, body.AllowedUsers)
	if err != nil {
		return err
	}
	web.JSONResponse(w, http.StatusCreated, createPostResponse{Post: p})
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *PostHandler) handleMultipartCreatePost(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	err := r.ParseMultipartForm(int64(web.MaxImageSize))
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

	p, err := h.serv.CreatePost(identity.ID, content, imageURL, privacy, groupID, allowedUsers)
	if err != nil {
		return err
	}
	web.JSONResponse(w, http.StatusCreated, createPostResponse{Post: p})
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return &models.ValidationError{Field: "auth", Message: "unauthorized"}
	}

	var content, imageURL string
	var postID int64

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(int64(web.MaxImageSize))
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

	c, err := h.serv.CreateComment(identity.ID, postID, content, imageURL)
	if err != nil {
		return err
	}
	web.JSONResponse(w, http.StatusCreated, c)
	return nil
}
