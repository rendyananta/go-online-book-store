package user

import (
	"context"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=authenticate.go -destination=auth_manager_mock_test.go -package user
type authManager interface {
	Token(ctx context.Context, userID string) (string, error)
}

type AuthenticatorUseCase struct {
	userRepo    userRepo
	authManager authManager
}

func NewAuthenticatorUseCase(userRepo userRepo, authManager authManager) (*AuthenticatorUseCase, error) {
	return &AuthenticatorUseCase{
		userRepo:    userRepo,
		authManager: authManager,
	}, nil
}

func (a AuthenticatorUseCase) Authenticate(ctx context.Context, param user.AuthenticateParam) (user.AuthenticateResult, error) {
	u, err := a.userRepo.FindByEmail(ctx, param.Email)
	if err != nil {
		return user.AuthenticateResult{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(param.Password))
	if err != nil {
		return user.AuthenticateResult{}, err
	}

	token, err := a.authManager.Token(ctx, u.ID)
	if err != nil {
		return user.AuthenticateResult{}, err
	}

	return user.AuthenticateResult{
		User: user.User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		},
		Token: token,
	}, nil
}
