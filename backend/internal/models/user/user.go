package user

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type contextKey string

const Key contextKey = "user_id"

// --------------------------------------------------------------------|

type User struct {
	ID          int64      `json:"id"`
	Email       string     `json:"email,omitempty"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Password    string     `json:"password,omitempty"`
	DOB         *time.Time `json:"dob,omitempty"`
	AvatarURL   string     `json:"-"`
	Nickname    string     `json:"nickname,omitempty"`
	AboutMe     string     `json:"about_me,omitempty"`
	ProfileType string     `json:"profile_type,omitempty"`
}

type UserData struct {
	ID        int64      `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	UUID      string     `json:"-"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// --------------------------------------------------------------------|

func (u *User) ValidateData() error {
	if err := u.isEmpty(); err != nil {
		return err
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return fmt.Errorf("incorrect email")
	}

	now := time.Now()

	if u.DOB.IsZero() {
		return fmt.Errorf("invalid birth date")
	}

	if u.DOB.After(now) {
		return fmt.Errorf("invalid birth date")
	}

	minDate := now.AddDate(-122, 0, 0)
	if u.DOB.Before(minDate) {
		return fmt.Errorf("invalid birth date")
	}
	return nil
}

// --------------------------------------------------------------------|

func (u *User) isEmpty() error {
	if len(strings.TrimSpace(u.Email)) == 0 {
		return fmt.Errorf("email is empty")
	}
	if len(strings.TrimSpace(u.FirstName)) == 0 {
		return fmt.Errorf("first name is empty")
	}
	if len(strings.TrimSpace(u.LastName)) == 0 {
		return fmt.Errorf("last name is empty")
	}
	return nil
}
