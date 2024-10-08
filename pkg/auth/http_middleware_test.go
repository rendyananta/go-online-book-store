package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type simpleErrorWriter struct {
}

func (s *simpleErrorWriter) Write(w http.ResponseWriter, r *http.Request, err error) {
	if !errors.Is(err, ErrUnauthenticated) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(err.Error()))
}

func TestMiddleware_Handle(t *testing.T) {
	type fields struct {
		auth      *Manager
		errWriter errorWriter
	}
	type args struct {
		next http.Handler
		req  *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		beforeTest     func(m *Middleware, args *args)
		want           string
		wantStatusCode int
	}{
		{
			name: "no token",
			fields: fields{
				auth:      &Manager{},
				errWriter: &simpleErrorWriter{},
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("success"))
				}),
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/", nil)
					return req
				}(),
			},
			beforeTest: func(m *Middleware, args *args) {

			},
			want:           ErrUnauthenticated.Error(),
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "invalid token header value",
			fields: fields{
				auth:      &Manager{},
				errWriter: &simpleErrorWriter{},
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("success"))
				}),
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/", nil)
					req.Header.Add(httpHeaderAuthKey, "Beareradasdas")
					return req
				}(),
			},
			beforeTest: func(m *Middleware, args *args) {

			},
			want:           ErrUnauthenticated.Error(),
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "token valid",
			fields: fields{
				auth: &Manager{
					config: Config{
						TokenLifetime: defaultTTL,
						CipherKeys:    []string{"0rMTKewMPeSGi6vi"},
					},
					cacheDriver: mockCacheDriver(),
					ciphers: []cipher.Block{
						func() cipher.Block {
							c, _ := aes.NewCipher([]byte("0rMTKewMPeSGi6vi"))
							return c
						}(),
					},
				},
				errWriter: &simpleErrorWriter{},
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("success"))
				}),
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/", nil)
					return req
				}(),
			},
			beforeTest: func(m *Middleware, args *args) {
				token, _ := m.auth.Token(context.Background(), fmt.Sprint(10))
				args.req.Header.Add(httpHeaderAuthKey, fmt.Sprintf("Bearer %s", token))
			},
			want:           "success",
			wantStatusCode: http.StatusOK,
		},
		{
			name: "token expired",
			fields: fields{
				auth: &Manager{
					config: Config{
						TokenLifetime: 10 * time.Millisecond,
						CipherKeys:    []string{"0rMTKewMPeSGi6vi"},
					},
					cacheDriver: mockCacheDriver(),
					ciphers: []cipher.Block{
						func() cipher.Block {
							c, _ := aes.NewCipher([]byte("0rMTKewMPeSGi6vi"))
							return c
						}(),
					},
				},
				errWriter: &simpleErrorWriter{},
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("success"))
				}),
				req: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/", nil)
					return req
				}(),
			},
			beforeTest: func(m *Middleware, args *args) {
				token, _ := m.auth.Token(context.Background(), fmt.Sprint(10))
				args.req.Header.Add(httpHeaderAuthKey, fmt.Sprintf("Bearer %s", token))

				time.Sleep(20 * time.Millisecond)
			},
			want:           ErrUnauthenticated.Error(),
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Middleware{
				auth:      tt.fields.auth,
				errWriter: tt.fields.errWriter,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(m, &tt.args)
			}

			got := m.Handle(tt.args.next)

			recorder := httptest.NewRecorder()

			got.ServeHTTP(recorder, tt.args.req)
			resp := recorder.Result()

			if !reflect.DeepEqual(resp.StatusCode, tt.wantStatusCode) {
				t.Errorf("Middleware.Handle() status code = %v, want %v", resp.StatusCode, tt.wantStatusCode)
			}

			body, _ := io.ReadAll(resp.Body)

			if !strings.EqualFold(string(body), tt.want) {
				t.Errorf("Middleware.Handle() body = %v, want %v", got, tt.want)
			}
		})
	}
}
