package models

import "time"

type PostAuthor struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type Post struct {
	ID           int64       `json:"id"`
	AuthorID     int64       `json:"author_id"`
	Author       *PostAuthor `json:"author,omitempty"`
	Content      string      `json:"content"`
	ImageURL     string      `json:"image_url,omitempty"`
	Privacy      string      `json:"privacy"`
	GroupID      int64       `json:"group_id,omitempty"`
	AllowedUsers []int64     `json:"allowed_users,omitempty"`
	Comments     []Comment   `json:"comments,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}

type Comment struct {
	ID        int64       `json:"id"`
	PostID    int64       `json:"post_id"`
	AuthorID  int64       `json:"author_id"`
	Author    *PostAuthor `json:"author,omitempty"`
	Content   string      `json:"content"`
	ImageURL  string      `json:"image_url,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type PostService interface {
	ListPosts(requesterID int64, groupID int64, limit, offset int) ([]Post, bool, error)
	CreatePost(authorID int64, content, imageURL, privacy string, groupID int64, allowedUsers []int64) (*Post, error)
	CreateComment(authorID int64, postID int64, content, imageURL string) (*Comment, error)
	GetComments(postID int64) ([]Comment, error)
}

type PostRepo interface {
	List(requesterID int64, groupID int64, limit, offset int) ([]Post, bool, error)
	Insert(authorID int64, content, imageURL, privacy string, groupID int64, allowedUsers []int64) (*Post, error)
	GetPosts(authorID int64, requesterID int64) ([]Post, error)
	GetGroupPosts(groupID int64, requesterID int64) ([]Post, error)
	InsertComment(authorID int64, postID int64, content, imageURL string) (*Comment, error)
	GetComments(postID int64) ([]Comment, error)
}

