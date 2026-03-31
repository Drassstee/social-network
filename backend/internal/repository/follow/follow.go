package follow

import "database/sql"

type FollowRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *FollowRepo {
	return &FollowRepo{db: db}
}
