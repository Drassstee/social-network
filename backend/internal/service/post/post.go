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

func (s *PostService) ListPosts(limit, offset int) ([]models.Post, bool, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(limit, offset)
}

func (s *PostService) CreatePost(authorID, content string) (*models.Post, error) {
	authorID = strings.TrimSpace(authorID)
	content = strings.TrimSpace(content)
	if authorID == "" {
		return nil, errors.New("author_id is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > maxPostContent {
		return nil, errors.New("content is too long")
	}
	return s.repo.Insert(authorID, content)
}
