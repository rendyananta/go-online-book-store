package order

import "errors"

const (
	defaultPaginationLength = 20
)

var (
	ErrNotFound = errors.New("not found")
)
