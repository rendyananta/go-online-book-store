package auth

import (
	"context"
	"net/http"
	"strings"
)

const (
	httpHeaderAuthKey = "Authorization"
	authTokenPrefix   = "Bearer "
)

type errorWriter interface {
	Write(w http.ResponseWriter, r *http.Request, err error)
}

type Middleware struct {
	auth      *Manager
	errWriter errorWriter
}

func NewMiddleware(authManager *Manager, httpErrWriter errorWriter) *Middleware {
	return &Middleware{
		auth:      authManager,
		errWriter: httpErrWriter,
	}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(httpHeaderAuthKey)

		if auth == "" {
			m.errWriter.Write(w, r, ErrUnauthenticated)
			return
		}

		after, ok := strings.CutPrefix(auth, authTokenPrefix)
		if !ok {
			m.errWriter.Write(w, r, ErrUnauthenticated)
			return
		}

		session, err := m.auth.User(r.Context(), after)
		if err != nil {
			m.errWriter.Write(w, r, ErrUnauthenticated)
			return
		}

		newCtx := context.WithValue(r.Context(), CtxKeyUserSession, &session)
		newReq := r.Clone(newCtx)

		next.ServeHTTP(w, newReq)
	})
}
