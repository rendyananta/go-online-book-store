package user

import "errors"

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrEmailIsNotRegistered   = errors.New("email is not registered")
	ErrNotFound               = errors.New("not found")
)
