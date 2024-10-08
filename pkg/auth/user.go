package auth

import (
	"context"
	"time"
)

func UserFromContext(ctx context.Context) (*UserSession, error) {
	user, ok := ctx.Value(CtxKeyUserSession).(*UserSession)
	if !ok {
		return nil, ErrUnauthenticated

	}

	if time.Now().UnixMilli() > user.ExpiredAt.UnixMilli() {
		return nil, ErrUnauthenticated
	}

	return user, nil
}
