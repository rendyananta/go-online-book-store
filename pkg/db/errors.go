package db

import "errors"

var (
	ErrConnectionUnregistered = errors.New("connection not registered")
)
