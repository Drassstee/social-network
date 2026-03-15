package models

import "time"

type User struct {
	ID        int64
	Email     string
	FirstName string
	LastName  string
	Password  string
	BirthDay  time.Time
	Avatar    string
	Nickname  string
	Bio       string
}
