package post

import (
	"database/sql"
	"time"

	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

type PostRepo struct {
	db *sql.DB
}

//--------------------------------------------------------------------------------------|

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

//--------------------------------------------------------------------------------------|

func (r *PostRepo) List(requesterID int64, groupID int64, limit, offset int) ([]models.Post, bool, error) {
	if limit < 1 {
		limit = 20
	}

	query := `
SELECT p.id, p.author_id, p.content, p.image_url, p.privacy, p.group_id, p.created_at, u.id, u.first_name, u.last_name, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE ( (? = 0 AND p.group_id IS NULL) OR (p.group_id = ?) )
  AND (
    (p.author_id = ?)
    OR (p.privacy = 'public' AND u.profile_type = 'public')
    OR ((p.privacy = 'public' OR p.privacy = 'almost_private') AND p.author_id IN (SELECT following_id FROM follows WHERE follower_id = ? AND status = 'accept'))
    OR (p.privacy = 'private' AND p.id IN (SELECT post_id FROM post_allowed_users WHERE user_id = ?))
  )
ORDER BY datetime(p.created_at) DESC
LIMIT ? OFFSET ?
`
	rows, err := r.db.Query(query, groupID, groupID, requesterID, requesterID, requesterID, limit+1, offset)

	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		var a models.PostAuthor
		var created string
		var img, priv, avt sql.NullString
		var grpID sql.NullInt64
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &img, &priv, &grpID, &created, &a.ID, &a.FirstName, &a.LastName, &avt); err != nil {
			return nil, false, err
		}
		a.AvatarURL = avt.String
		p.Author = &a
		p.ImageURL = img.String
		p.Privacy = priv.String
		if grpID.Valid {
			p.GroupID = grpID.Int64
		}
		p.CreatedAt, _ = parseSQLiteTime(created)

		posts = append(posts, p)
	}
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}
	return posts, hasMore, nil
}

//--------------------------------------------------------------------------------------|

func (r *PostRepo) Insert(authorID int64, content, imageURL, privacy string, groupID int64, allowedUsers []int64) (*models.Post, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var gID *int64
	if groupID > 0 {
		gID = &groupID
	}

	res, err := tx.Exec(`INSERT INTO posts (author_id, content, image_url, privacy, group_id) VALUES (?, ?, ?, ?, ?)`,
		authorID, content, imageURL, privacy, gID)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	if privacy == "private" && groupID == 0 {
		for _, userID := range allowedUsers {
			_, err = tx.Exec(`INSERT INTO post_allowed_users (post_id, user_id) VALUES (?, ?)`, id, userID)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Fetch the inserted post
	row := r.db.QueryRow(`
SELECT p.id, p.author_id, p.content, p.image_url, p.privacy, p.group_id, p.created_at, u.id, u.first_name, u.last_name, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.id = ?
`, id)
	var p models.Post
	var a models.PostAuthor
	var created string
	var img, priv, avt sql.NullString
	var grpID sql.NullInt64
	if err := row.Scan(&p.ID, &p.AuthorID, &p.Content, &img, &priv, &grpID, &created, &a.ID, &a.FirstName, &a.LastName, &avt); err != nil {
		return nil, err
	}
	a.AvatarURL = avt.String
	p.Author = &a
	p.ImageURL = img.String
	p.Privacy = priv.String
	if grpID.Valid {
		p.GroupID = grpID.Int64
	}
	p.CreatedAt, _ = parseSQLiteTime(created)
	return &p, nil
}

//--------------------------------------------------------------------------------------|

func (r *PostRepo) GetPosts(authorID int64, requesterID int64) ([]models.Post, error) {
	rows, err := r.db.Query(`
SELECT p.id, p.author_id, p.content, p.image_url, p.privacy, p.group_id, p.created_at, u.id, u.first_name, u.last_name, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.author_id = ? AND p.group_id IS NULL
  AND ( (p.author_id = ?)
     OR (p.privacy = 'public' AND u.profile_type = 'public')
     OR ((p.privacy = 'public' OR p.privacy = 'almost_private') AND p.author_id IN (SELECT following_id FROM follows WHERE follower_id = ? AND status = 'accept'))
     OR (p.privacy = 'private' AND p.id IN (SELECT post_id FROM post_allowed_users WHERE user_id = ?))
  )
ORDER BY datetime(p.created_at) DESC
`, authorID, requesterID, requesterID, requesterID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		var a models.PostAuthor
		var created string
		var img, priv, avt sql.NullString
		var grpID sql.NullInt64
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &img, &priv, &grpID, &created, &a.ID, &a.FirstName, &a.LastName, &avt); err != nil {
			return nil, err
		}
		a.AvatarURL = avt.String
		p.Author = &a
		p.ImageURL = img.String
		p.Privacy = priv.String
		if grpID.Valid {
			p.GroupID = grpID.Int64
		}
		p.CreatedAt, _ = parseSQLiteTime(created)
		posts = append(posts, p)
	}
	return posts, nil
}

//--------------------------------------------------------------------------------------|

func (r *PostRepo) GetGroupPosts(groupID int64, requesterID int64) ([]models.Post, error) {
	// Check if user is member of the group
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, requesterID).Scan(&count)
	if err != nil || count == 0 {
		return nil, nil // Or error? Instructions say only displayed to members.
	}

	rows, err := r.db.Query(`
SELECT p.id, p.author_id, p.content, p.image_url, p.privacy, p.group_id, p.created_at, u.id, u.first_name, u.last_name, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.group_id = ?
ORDER BY datetime(p.created_at) DESC
`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		var a models.PostAuthor
		var created string
		var img, priv, avt sql.NullString
		var gID sql.NullInt64
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &img, &priv, &gID, &created, &a.ID, &a.FirstName, &a.LastName, &avt); err != nil {
			return nil, err
		}
		a.AvatarURL = avt.String
		p.Author = &a
		p.ImageURL = img.String
		p.Privacy = priv.String
		if gID.Valid {
			p.GroupID = gID.Int64
		}
		p.CreatedAt, _ = parseSQLiteTime(created)

		// Load comments
		comments, _ := r.GetComments(p.ID)
		p.Comments = comments

		posts = append(posts, p)
	}
	return posts, nil
}

//--------------------------------------------------------------------------------------|

func (r *PostRepo) InsertComment(authorID int64, postID int64, content, imageURL string) (*models.Comment, error) {
	res, err := r.db.Exec(`INSERT INTO comments (post_id, author_id, content, image_url) VALUES (?, ?, ?, ?)`,
		postID, authorID, content, imageURL)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	row := r.db.QueryRow(`
SELECT c.id, c.post_id, c.author_id, c.content, c.image_url, c.created_at, u.id, u.first_name, u.last_name, u.avatar_url
FROM comments c
JOIN users u ON c.author_id = u.id
WHERE c.id = ?
`, id)
	var c models.Comment
	var a models.PostAuthor
	var created string
	var img, avt sql.NullString
	if err := row.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &img, &created, &a.ID, &a.FirstName, &a.LastName, &avt); err != nil {
		return nil, err
	}
	a.AvatarURL = avt.String
	c.Author = &a
	c.ImageURL = img.String
	c.CreatedAt, _ = parseSQLiteTime(created)
	return &c, nil
}

//--------------------------------------------------------------------------------------|

func (r *PostRepo) GetComments(postID int64) ([]models.Comment, error) {
	rows, err := r.db.Query(`
SELECT c.id, c.post_id, c.author_id, c.content, c.image_url, c.created_at, u.id, u.first_name, u.last_name, u.avatar_url
FROM comments c
JOIN users u ON c.author_id = u.id
WHERE c.post_id = ?
ORDER BY datetime(c.created_at) ASC
`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		var a models.PostAuthor
		var created string
		var img, avt sql.NullString
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &img, &created, &a.ID, &a.FirstName, &a.LastName, &avt); err != nil {
			continue
		}
		a.AvatarURL = avt.String
		c.Author = &a
		c.ImageURL = img.String
		c.CreatedAt, _ = parseSQLiteTime(created)
		comments = append(comments, c)
	}
	return comments, nil
}

//--------------------------------------------------------------------------------------|

func (r *PostRepo) GetCommentsBatch(postIDs []int64) (map[int64][]models.Comment, error) {
	if len(postIDs) == 0 {
		return make(map[int64][]models.Comment), nil
	}

	placeholders := ""
	args := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	query := `
SELECT c.id, c.post_id, c.author_id, c.content, c.image_url, c.created_at, u.id, u.first_name, u.last_name, u.avatar_url
FROM comments c
JOIN users u ON c.author_id = u.id
WHERE c.post_id IN (` + placeholders + `)
ORDER BY datetime(c.created_at) ASC
`
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[int64][]models.Comment)
	for rows.Next() {
		var c models.Comment
		var a models.PostAuthor
		var created string
		var img, avt sql.NullString
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &img, &created, &a.ID, &a.FirstName, &a.LastName, &avt); err != nil {
			continue
		}
		a.AvatarURL = avt.String
		c.Author = &a
		c.ImageURL = img.String
		c.CreatedAt, _ = parseSQLiteTime(created)
		results[c.PostID] = append(results[c.PostID], c)
	}
	return results, nil
}

//--------------------------------------------------------------------------------------|

func parseSQLiteTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.UTC)
}
