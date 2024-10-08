package auth

import "time"

const (
	defaultUserType = ""
	defaultTTL      = 60 * time.Minute
)

type CtxKey string

const CtxKeyUserSession CtxKey = "user_session"
