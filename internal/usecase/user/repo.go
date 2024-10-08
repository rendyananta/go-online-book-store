package user

import (
	"context"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
)

//go:generate mockgen -source=repo.go -destination=repo_mock_test.go -package user
type userRepo interface {
	FindByEmail(ctx context.Context, email string) (user.User, error)
	Create(ctx context.Context, param user.User) (user.User, error)
}
