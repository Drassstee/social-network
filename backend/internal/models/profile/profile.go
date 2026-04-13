package profile

import (
"social-network/internal/models"
"social-network/internal/models/user"
)

type Profile struct {
User      *user.User      `json:"user"`
Followers []user.UserData `json:"followers"`
Following []user.UserData `json:"following"`
Posts     []models.Post   `json:"posts"`
}
