package profile

import (
	"social-network/internal/models/follow"
	"social-network/internal/models/post"
	"social-network/internal/models/user"
)

type Profile struct {
	User      *user.User
	Posts     []post.Post
	Followers []follow.Follow
	Following []follow.Follow
}
