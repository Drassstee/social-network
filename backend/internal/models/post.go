package models

import "time"

type PostAuthor struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Post struct {
	ID        int64       `json:"id"`
	AuthorID  int64       `json:"author_id"`
	Author    *PostAuthor `json:"author,omitempty"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
}

type PostService interface {
	ListPosts(limit, offset int) ([]Post, bool, error)
	CreatePost(authorID int64, content string) (*Post, error)
}

type PostRepo interface {
	List(limit, offset int) ([]Post, bool, error)
	Insert(authorID int64, content string) (*Post, error)
	GetPosts(authorID int64) ([]Post, error)
}
