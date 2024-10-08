package auth

import "errors"

var (
	ErrCipherKeysIsEmpty  error = errors.New("cipher keys config is empty")
	ErrUnauthenticated    error = errors.New("unauthenticated")
	ErrTokenExpired       error = errors.New("token expired")
	ErrSessionKeyNotFound error = errors.New("session key not found")
	ErrInvalidTokenSize   error = errors.New("invalid token size")
)
