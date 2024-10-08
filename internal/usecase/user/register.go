package user

import (
	"context"
	"errors"

	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	repo "github.com/rendyananta/example-online-book-store/internal/repo/user"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUseCase struct {
	userRepo userRepo
}

func NewRegisterUseCase(userRepo userRepo) (*RegisterUseCase, error) {
	return &RegisterUseCase{
		userRepo: userRepo,
	}, nil
}

func (r RegisterUseCase) Register(ctx context.Context, param user.RegisterParam) (user.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)

	u, err := r.userRepo.Create(ctx, user.User{
		Name:     param.Name,
		Email:    param.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		return u, err
	}

	return u, nil
}

func (r RegisterUseCase) EmailRegistered(ctx context.Context, email string) bool {
	u, err := r.userRepo.FindByEmail(ctx, email)
	if err != nil && errors.Is(err, repo.ErrEmailIsNotRegistered) {
		return false
	}

	if u.ID == "" {
		return false
	}

	return true
}
