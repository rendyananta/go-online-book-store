package cache

import "errors"

var (
	ErrDriverUnregistered = errors.New("driver not registered")
	ErrNotFound           = errors.New("not found")
)
