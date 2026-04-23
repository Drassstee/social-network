package post

import (
	"errors"
	"strings"

	"social-network/internal/models"
)

const (
	maxPostContent = 16000
	defaultLimit   = 20
	maxLimit       = 100
)

type PostService struct {
	repo models.PostRepo
}

func NewPostService(repo models.PostRepo) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) ListPosts(requesterID int64, groupID int64, limit, offset int) ([]models.Post, bool, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	if groupID > 0 {
		posts, err := s.repo.GetGroupPosts(groupID, requesterID)
		if err != nil {
			return nil, false, err
		}
		// Basic pagination for memory-cached group posts (or I could implement it in repo)
		// For now simple return as is or implement paging.
		return posts, false, nil
	}
	return s.repo.List(requesterID, 0, limit, offset)
}

func (s *PostService) CreatePost(authorID int64, content, imageURL, privacy string, groupID int64, allowedUsers []int64) (*models.Post, error) {
	content = strings.TrimSpace(content)
	if authorID == 0 {
		return nil, errors.New("author_id is required")
	}
	if content == "" && imageURL == "" {
		return nil, errors.New("content or image is required")
	}
	if len(content) > maxPostContent {
		return nil, errors.New("content is too long")
	}
	if privacy == "" {
		privacy = "public"
	}
	return s.repo.Insert(authorID, content, imageURL, privacy, groupID, allowedUsers)
}

func (s *PostService) CreateComment(authorID int64, postID int64, content, imageURL string) (*models.Comment, error) {
	content = strings.TrimSpace(content)
	if authorID == 0 || postID == 0 {
		return nil, errors.New("author_id and post_id are required")
	}
	if content == "" && imageURL == "" {
		return nil, errors.New("content or image is required")
	}
	return s.repo.InsertComment(authorID, postID, content, imageURL)
}

func (s *PostService) GetComments(postID int64) ([]models.Comment, error) {
	return s.repo.GetComments(postID)
}

