package app

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal")
	ErrNotFound     = errors.New("not found")
)
