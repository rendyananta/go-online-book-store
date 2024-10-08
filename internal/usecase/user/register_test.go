package user

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	userrp "github.com/rendyananta/example-online-book-store/internal/repo/user"
	"reflect"
	"testing"
)

func TestNewRegisterUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := NewMockuserRepo(ctrl)

	type args struct {
		userRepo userRepo
	}
	tests := []struct {
		name    string
		args    args
		want    *RegisterUseCase
		wantErr bool
	}{
		{
			name: "init new use case",
			args: args{
				userRepo: userRepoMock,
			},
			want: &RegisterUseCase{
				userRepo: userRepoMock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRegisterUseCase(tt.args.userRepo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegisterUseCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegisterUseCase() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterUseCase_EmailRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := NewMockuserRepo(ctrl)

	type fields struct {
		userRepo userRepo
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(uc *RegisterUseCase)
		want       bool
	}{
		{
			name: "email has been registered",
			fields: fields{
				userRepo: userRepoMock,
			},
			args: args{
				ctx:   context.Background(),
				email: "example@email.com",
			},
			beforeTest: func(uc *RegisterUseCase) {
				userRepoMock.EXPECT().FindByEmail(context.Background(), "example@email.com").
					Return(user.User{
						ID:    "1234",
						Name:  "Example email",
						Email: "example@email.com",
					}, nil)
			},
			want: true,
		},
		{
			name: "email is not registered",
			fields: fields{
				userRepo: userRepoMock,
			},
			args: args{
				ctx:   context.Background(),
				email: "example@email.com",
			},
			beforeTest: func(uc *RegisterUseCase) {
				userRepoMock.EXPECT().FindByEmail(context.Background(), "example@email.com").
					Return(user.User{}, userrp.ErrEmailIsNotRegistered)
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RegisterUseCase{
				userRepo: tt.fields.userRepo,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(&r)
			}

			if got := r.EmailRegistered(tt.args.ctx, tt.args.email); got != tt.want {
				t.Errorf("EmailRegistered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := NewMockuserRepo(ctrl)

	type fields struct {
		userRepo userRepo
	}
	type args struct {
		ctx   context.Context
		param user.RegisterParam
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(uc *RegisterUseCase)
		want       user.User
		wantErr    bool
	}{
		{
			name: "can handle registration",
			fields: fields{
				userRepo: userRepoMock,
			},
			args: args{
				ctx: context.Background(),
				param: user.RegisterParam{
					Name:     "Example User",
					Email:    "example@email.com",
					Password: "12345678",
				},
			},
			beforeTest: func(uc *RegisterUseCase) {
				userRepoMock.EXPECT().Create(context.Background(), userMatcher{input: user.User{
					Name:  "Example User",
					Email: "example@email.com",
				}}).Return(user.User{
					ID:       "1234",
					Name:     "Example User",
					Email:    "example@email.com",
					Password: "hashed-password",
				}, nil)
			},
			want: user.User{
				ID:       "1234",
				Name:     "Example User",
				Email:    "example@email.com",
				Password: "hashed-password",
			},
			wantErr: false,
		},
		{
			name: "can handle error in registration",
			fields: fields{
				userRepo: userRepoMock,
			},
			args: args{
				ctx: context.Background(),
				param: user.RegisterParam{
					Name:     "Example User",
					Email:    "example@email.com",
					Password: "12345678",
				},
			},
			beforeTest: func(uc *RegisterUseCase) {
				userRepoMock.EXPECT().Create(context.Background(), userMatcher{input: user.User{
					Name:  "Example User",
					Email: "example@email.com",
				}}).Return(user.User{}, sql.ErrConnDone)
			},
			want:    user.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RegisterUseCase{
				userRepo: tt.fields.userRepo,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(&r)
			}

			got, err := r.Register(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}
