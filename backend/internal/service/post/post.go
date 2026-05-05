package post

import (
	"errors"
	"strings"

	"social-network/internal/models"
	"social-network/internal/utils"
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
	posts, hasMore, err := s.repo.List(requesterID, 0, limit, offset)
	if err != nil {
		return nil, false, err
	}
	for i := range posts {
		if posts[i].Author != nil {
			posts[i].Author.AvatarURL = utils.FormatAvatarURL(posts[i].Author.AvatarURL)
		}
		for j := range posts[i].Comments {
			if posts[i].Comments[j].Author != nil {
				posts[i].Comments[j].Author.AvatarURL = utils.FormatAvatarURL(posts[i].Comments[j].Author.AvatarURL)
			}
		}
	}
	return posts, hasMore, nil
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
	p, err := s.repo.Insert(authorID, content, imageURL, privacy, groupID, allowedUsers)
	if err == nil && p.Author != nil {
		p.Author.AvatarURL = utils.FormatAvatarURL(p.Author.AvatarURL)
	}
	return p, err
}

func (s *PostService) CreateComment(authorID int64, postID int64, content, imageURL string) (*models.Comment, error) {
	content = strings.TrimSpace(content)
	if authorID == 0 || postID == 0 {
		return nil, errors.New("author_id and post_id are required")
	}
	if content == "" && imageURL == "" {
		return nil, errors.New("content or image is required")
	}
	c, err := s.repo.InsertComment(authorID, postID, content, imageURL)
	if err == nil && c.Author != nil {
		c.Author.AvatarURL = utils.FormatAvatarURL(c.Author.AvatarURL)
	}
	return c, err
}

func (s *PostService) GetComments(postID int64) ([]models.Comment, error) {
	comments, err := s.repo.GetComments(postID)
	if err == nil {
		for i := range comments {
			if comments[i].Author != nil {
				comments[i].Author.AvatarURL = utils.FormatAvatarURL(comments[i].Author.AvatarURL)
			}
		}
	}
	return comments, err
}

