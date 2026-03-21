package models

import "time"

type Post struct {
	ID        int64     `json:"id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type PostService interface {
	ListPosts(limit, offset int) ([]Post, bool, error)
	CreatePost(authorID, content string) (*Post, error)
}

type PostRepo interface {
	List(limit, offset int) ([]Post, bool, error)
	Insert(authorID, content string) (*Post, error)
}
