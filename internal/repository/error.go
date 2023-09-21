package repository

import "errors"

var (
	ErrUnique     = errors.New("not unique")
	ErrValidation = errors.New("invalid")
	ErrNotFound   = errors.New("not found")
)
