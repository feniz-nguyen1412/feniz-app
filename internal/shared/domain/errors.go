package domain

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidID        = errors.New("invalid id")
	ErrDuplicateEmail   = errors.New("duplicate email")
)