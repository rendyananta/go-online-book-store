package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type arrayCacheDriver struct {
	array map[string][]byte
	opt   mockFuncOpt
}

func (a *arrayCacheDriver) Get(ctx context.Context, key string) ([]byte, error) {
	if a.opt.getFunc != nil {
		return a.opt.getFunc()
	}

	val, ok := a.array[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	return val, nil
}

func (a *arrayCacheDriver) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	if a.opt.setFunc != nil {
		return a.opt.setFunc()
	}

	a.array[key] = val
	return nil
}

func (a *arrayCacheDriver) Del(ctx context.Context, key string) error {
	if a.opt.delFunc != nil {
		return a.opt.delFunc()
	}

	delete(a.array, key)
	return nil
}

type mockFuncOpt struct {
	getFunc func() ([]byte, error)
	setFunc func() error
	delFunc func() error
}

func mockCacheDriver() *arrayCacheDriver {
	return &arrayCacheDriver{
		array: make(map[string][]byte, 0),
	}
}

func mockCacheDriverWithOpt(opt mockFuncOpt) *arrayCacheDriver {
	return &arrayCacheDriver{
		array: make(map[string][]byte, 0),
		opt:   opt,
	}
}

func TestNewAuthManager(t *testing.T) {
	type args struct {
		conf        Config
		cacheDriver cacheDriver
	}
	tests := []struct {
		name    string
		args    args
		want    *Manager
		wantErr bool
	}{
		{
			name: "can init",
			args: args{
				conf: Config{
					CipherKeys: []string{"0rMTKewMPeSGi6vi", "Vo6g1ixi33zxc2Kb"},
				},
				cacheDriver: mockCacheDriver(),
			},
			want: &Manager{
				config: Config{
					TokenLifetime: defaultTTL,
					CipherKeys:    []string{"0rMTKewMPeSGi6vi", "Vo6g1ixi33zxc2Kb"},
				},
				ciphers: []cipher.Block{
					func() cipher.Block {
						c, _ := aes.NewCipher([]byte("Vo6g1ixi33zxc2Kb"))
						return c
					}(),
					func() cipher.Block {
						c, _ := aes.NewCipher([]byte("0rMTKewMPeSGi6vi"))
						return c
					}(),
				},
				cacheDriver: mockCacheDriver(),
			},
			wantErr: false,
		},
		{
			name: "can handle empty cipher",
			args: args{
				conf:        Config{},
				cacheDriver: mockCacheDriver(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "can handle invalid cipher length",
			args: args{
				conf: Config{
					CipherKeys: []string{"key1"},
				},
				cacheDriver: mockCacheDriver(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAuthManager(tt.args.conf, tt.args.cacheDriver)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuthManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthManager_Token(t *testing.T) {
	type fields struct {
		config      Config
		cacheDriver cacheDriver
		ciphers     []cipher.Block
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantStringVal bool
		wantErr       bool
	}{
		{
			name: "can generate token",
			fields: fields{
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
			args: args{
				ctx:    context.Background(),
				userID: fmt.Sprint(10),
			},
			wantStringVal: true,
			wantErr:       false,
		},
		{
			name: "can handle fail generate token",
			fields: fields{
				config: Config{
					TokenLifetime: defaultTTL,
					CipherKeys:    []string{"0rMTKewMPeSGi6vi"},
				},
				cacheDriver: mockCacheDriverWithOpt(mockFuncOpt{
					setFunc: func() error {
						return errors.New("failed to set session")
					},
				}),
				ciphers: []cipher.Block{
					func() cipher.Block {
						c, _ := aes.NewCipher([]byte("0rMTKewMPeSGi6vi"))
						return c
					}(),
				},
			},
			args: args{
				ctx:    context.Background(),
				userID: fmt.Sprint(10),
			},
			wantStringVal: false,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				config:      tt.fields.config,
				cacheDriver: tt.fields.cacheDriver,
				ciphers:     tt.fields.ciphers,
			}
			got, err := a.Token(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != "") != tt.wantStringVal {
				t.Errorf("Manager.Token() = %v, want %v", got, tt.wantStringVal)
			}
		})
	}
}

func TestAuthManager_User(t *testing.T) {
	type fields struct {
		config      Config
		cacheDriver cacheDriver
		ciphers     []cipher.Block
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(a *Manager, args *args)
		want       UserSession
		wantErr    bool
	}{
		{
			name: "can resolve user from token",
			fields: fields{
				config: Config{
					TokenLifetime: defaultTTL,
					CipherKeys:    []string{"0rMTKewMPeSGi6vi"},
				},
				cacheDriver: mockCacheDriverWithOpt(mockFuncOpt{
					getFunc: func() ([]byte, error) {
						val, _ := json.Marshal(UserSession{
							ID:        fmt.Sprint(10),
							ExpiredAt: time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
						})

						return val, nil
					},
				}),
				ciphers: []cipher.Block{
					func() cipher.Block {
						c, _ := aes.NewCipher([]byte("0rMTKewMPeSGi6vi"))
						return c
					}(),
				},
			},
			args: args{
				ctx:   context.Background(),
				token: "",
			},
			beforeTest: func(a *Manager, args *args) {
				token, _ := a.Token(context.Background(), fmt.Sprint(10))
				args.token = token
			},
			want: UserSession{
				ID:        fmt.Sprint(10),
				ExpiredAt: time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "can block expired token",
			fields: fields{
				config: Config{
					TokenLifetime: defaultTTL,
					CipherKeys:    []string{"0rMTKewMPeSGi6vi"},
				},
				cacheDriver: mockCacheDriverWithOpt(mockFuncOpt{
					getFunc: func() ([]byte, error) {
						val, _ := json.Marshal(UserSession{
							ID:        fmt.Sprint(10),
							ExpiredAt: time.Now().Add(-2 * time.Hour),
						})

						return val, nil
					},
				}),
				ciphers: []cipher.Block{
					func() cipher.Block {
						c, _ := aes.NewCipher([]byte("0rMTKewMPeSGi6vi"))
						return c
					}(),
				},
			},
			args: args{
				ctx:   context.Background(),
				token: "",
			},
			beforeTest: func(a *Manager, args *args) {
				token, _ := a.Token(context.Background(), fmt.Sprint(10))
				args.token = token
			},
			want:    UserSession{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				config:      tt.fields.config,
				cacheDriver: tt.fields.cacheDriver,
				ciphers:     tt.fields.ciphers,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(a, &tt.args)
			}

			got, err := a.User(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.User() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.User() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthManager_sessionKeyFor(t *testing.T) {
	type fields struct {
		config      Config
		cacheDriver cacheDriver
		ciphers     []cipher.Block
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(a *Manager, args *args)
		want       string
		wantErr    bool
	}{
		{
			name: "can decrypt token key",
			fields: fields{
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
			args: args{
				ctx:   context.Background(),
				token: "",
			},
			beforeTest: func(a *Manager, args *args) {
				token, _ := a.Token(context.Background(), fmt.Sprint(10))
				args.token = token
			},
			want:    "auth:_10_",
			wantErr: false,
		},
		{
			name: "can handle invalid token",
			fields: fields{
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
			args: args{
				ctx:   context.Background(),
				token: "Bdkaoq0uZ0POHFKSs5Nf5jqbob13jstdXOi8wMpQQ0rqtESf+j7J5EUB+678vjprM9MhNcsYtWfIaa", // invalid token
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				config:      tt.fields.config,
				cacheDriver: tt.fields.cacheDriver,
				ciphers:     tt.fields.ciphers,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(a, &tt.args)
			}

			got, err := a.sessionKeyFor(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.sessionKeyFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(got, tt.want) {
				t.Errorf("Manager.sessionKeyFor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthManager_Revoke(t *testing.T) {
	type fields struct {
		config      Config
		cacheDriver cacheDriver
		ciphers     []cipher.Block
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(a *Manager, args *args)
		wantErr    bool
	}{
		{
			name: "can revoke active token",
			fields: fields{
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
			args: args{
				ctx:   context.Background(),
				token: "",
			},
			beforeTest: func(a *Manager, args *args) {
				token, _ := a.Token(context.Background(), fmt.Sprint(10))
				args.token = token
			},
			wantErr: false,
		},
		{
			name: "can revoke non existing token",
			fields: fields{
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
			args: args{
				ctx:   context.Background(),
				token: "Bdkaoq0uZ0POHFKSs5Nf5jqbob13jstdXOi8wMpQQ0rqtESf+j7J5EUB+678vjprM9MhNcsYtWfIaa", // invalid token
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				config:      tt.fields.config,
				cacheDriver: tt.fields.cacheDriver,
				ciphers:     tt.fields.ciphers,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(a, &tt.args)
			}

			if err := a.Revoke(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Revoke() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
