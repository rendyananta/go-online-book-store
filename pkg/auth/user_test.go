package auth

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestUserFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *UserSession
		wantErr bool
	}{
		{
			name: "authenticated",
			args: args{
				ctx: context.WithValue(context.Background(), CtxKeyUserSession, &UserSession{
					ID:        fmt.Sprint(1),
					ExpiredAt: time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
				}),
			},
			want: &UserSession{
				ID:        fmt.Sprint(1),
				ExpiredAt: time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "expired",
			args: args{
				ctx: context.WithValue(context.Background(), CtxKeyUserSession, &UserSession{
					ID:        fmt.Sprint(1),
					ExpiredAt: time.Now().Add(-10 * time.Second),
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ctx doesn't have user session key",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserFromContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
