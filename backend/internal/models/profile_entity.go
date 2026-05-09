package models

import "time"

//--------------------------------------------------------------------------------------|

type Profile struct {
	User         UserData   `json:"user"`
	Followers    []UserData `json:"followers"`
	Following    []UserData `json:"following"`
	Posts        []Post     `json:"posts"`
	IsMe         bool       `json:"is_me"`
	FollowStatus string     `json:"follow_status"`
	Privacy      string     `json:"privacy"`
	DOB          *time.Time `json:"dob"`
	AboutMe      string     `json:"about_me"`
	Email        string     `json:"email"`
}

//--------------------------------------------------------------------------------------|

// UserProfile is used for the profile view
type UserProfile struct {
	ID           int64      `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Nickname     string     `json:"nickname"`
	AvatarURL    string     `json:"avatar_url"`
	Privacy      string     `json:"privacy"`
	DOB          *time.Time `json:"dob"`
	AboutMe      string     `json:"about_me"`
	IsMe         bool       `json:"is_me"`
	FollowStatus string     `json:"follow_status"`
}
