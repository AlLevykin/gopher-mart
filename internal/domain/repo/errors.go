package repo

import "errors"

var (
	ErrUserExists     = errors.New("user already exists")
	ErrUserValidation = errors.New("user validation error")
)
