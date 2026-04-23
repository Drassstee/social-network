package post

import (
	"database/sql"
	"time"

	"social-network/internal/models"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) List(limit, offset int) ([]models.Post, bool, error) {
	if limit < 1 {
		limit = 20
	}
	rows, err := r.db.Query(`
SELECT p.id, p.author_id, p.content, p.created_at, u.id, u.first_name, u.last_name
FROM posts p
JOIN users u ON p.author_id = u.id
ORDER BY datetime(p.created_at) DESC
LIMIT ? OFFSET ?
`, limit+1, offset)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var a models.PostAuthor
		var created string
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &created, &a.ID, &a.FirstName, &a.LastName); err != nil {
			return nil, false, err
		}
		p.Author = &a
		p.CreatedAt, err = parseSQLiteTime(created)
		if err != nil {
			return nil, false, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}
	return posts, hasMore, nil
}

func (r *PostRepo) Insert(authorID int64, content string) (*models.Post, error) {
	res, err := r.db.Exec(`INSERT INTO posts (author_id, content) VALUES (?, ?)`, authorID, content)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(`
SELECT p.id, p.author_id, p.content, p.created_at, u.id, u.first_name, u.last_name
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.id = ?
`, id)
	var p models.Post
	var a models.PostAuthor
	var created string
	if err := row.Scan(&p.ID, &p.AuthorID, &p.Content, &created, &a.ID, &a.FirstName, &a.LastName); err != nil {
		return nil, err
	}
	p.Author = &a
	p.CreatedAt, err = parseSQLiteTime(created)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepo) GetPosts(authorID int64) ([]models.Post, error) {
	rows, err := r.db.Query(`
SELECT p.id, p.author_id, p.content, p.created_at, u.id, u.first_name, u.last_name
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.author_id = ?
`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var a models.PostAuthor
		var created string
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &created, &a.ID, &a.FirstName, &a.LastName); err != nil {
			return nil, err
		}
		p.Author = &a
		p.CreatedAt, _ = parseSQLiteTime(created)
		posts = append(posts, p)
	}
	return posts, nil
}

func parseSQLiteTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.UTC)
}
