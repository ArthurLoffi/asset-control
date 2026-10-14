package exception

import "errors"

var (
	ErrNotFound   = errors.New("record not found")
	ErrConflict   = errors.New("record already exists")
	ErrValidation = errors.New("invalid input")
)
