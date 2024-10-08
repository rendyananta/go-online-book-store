package user

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	"reflect"
	"testing"
)

func TestAuthenticatorUseCase_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := NewMockuserRepo(ctrl)
	authManagerMock := NewMockauthManager(ctrl)

	type fields struct {
		userRepo    userRepo
		authManager authManager
	}
	type args struct {
		ctx   context.Context
		param user.AuthenticateParam
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       user.AuthenticateResult
		wantErr    bool
	}{
		{
			name: "can handle authentication",
			fields: fields{
				userRepo:    userRepoMock,
				authManager: authManagerMock,
			},
			args: args{
				ctx: context.Background(),
				param: user.AuthenticateParam{
					Email:    "user@example.com",
					Password: "123123",
				},
			},
			beforeTest: func() {
				userRepoMock.EXPECT().FindByEmail(context.Background(), "user@example.com").Return(user.User{
					ID:       "1",
					Name:     "User",
					Email:    "user@example.com",
					Password: "$2a$10$kzKHrJg9yufBEw3bpbUU8uoEtjAN3sREqWNR/b8eyX3s./1xSaAkq",
				}, nil)

				authManagerMock.EXPECT().Token(context.Background(), "1").
					Return("token-example", nil)
			},
			want: user.AuthenticateResult{
				User: user.User{
					ID:    "1",
					Name:  "User",
					Email: "user@example.com",
				},
				Token: "token-example",
			},
			wantErr: false,
		},
		{
			name: "can handle invalid authentication",
			fields: fields{
				userRepo:    userRepoMock,
				authManager: authManagerMock,
			},
			args: args{
				ctx: context.Background(),
				param: user.AuthenticateParam{
					Email:    "user@example.com",
					Password: "1231234",
				},
			},
			beforeTest: func() {
				userRepoMock.EXPECT().FindByEmail(context.Background(), "user@example.com").Return(user.User{
					ID:       "1",
					Name:     "User",
					Email:    "user@example.com",
					Password: "$2a$10$kzKHrJg9yufBEw3bpbUU8uoEtjAN3sREqWNR/b8eyX3s./1xSaAkq",
				}, nil)
			},
			want:    user.AuthenticateResult{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AuthenticatorUseCase{
				userRepo:    tt.fields.userRepo,
				authManager: tt.fields.authManager,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := a.Authenticate(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authenticate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAuthenticatorUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := NewMockuserRepo(ctrl)
	authManagerMock := NewMockauthManager(ctrl)

	type args struct {
		userRepo    userRepo
		authManager authManager
	}
	tests := []struct {
		name    string
		args    args
		want    *AuthenticatorUseCase
		wantErr bool
	}{
		{
			name: "can construct authenticator with nil",
			args: args{
				userRepo:    userRepoMock,
				authManager: authManagerMock,
			},
			want: &AuthenticatorUseCase{
				userRepo:    userRepoMock,
				authManager: authManagerMock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAuthenticatorUseCase(tt.args.userRepo, tt.args.authManager)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuthenticatorUseCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthenticatorUseCase() got = %v, want %v", got, tt.want)
			}
		})
	}
}
