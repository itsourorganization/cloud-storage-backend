package services

import "errors"

var (
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenInvalid   = errors.New("token invalid")
	ErrNotFound       = errors.New("not found")
	ErrIncorrectInput = errors.New("incorrect input")
	ErrUnique         = errors.New("not unique")
)
