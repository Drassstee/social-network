package follow

type Follow struct {
	FollowerID  int64  `json:"follower_id,omitempty"`
	FollowingID int64  `json:"following_id,omitempty"`
	Status      string `json:"status,omitempty"`
}
