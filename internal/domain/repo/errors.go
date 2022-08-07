package repo

import "errors"

var (
	ErrUserExists = errors.New("user already exists")
	ErrValidation = errors.New("validation error")
)
