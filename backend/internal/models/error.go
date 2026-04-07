package models

import "errors"

var (
	ErrInvalidData = errors.New("invalid data")       //400
	ErrUserPrivate = errors.New("profile is private") // 403
	ErrNotFound    = errors.New("not found")          // 404
	ErrConflict    = errors.New("conflict")           // 409
)
