package main

import (
	"github.com/rendyananta/example-online-book-store/internal/presenter/http/book"
	"github.com/rendyananta/example-online-book-store/internal/presenter/http/order"
	"github.com/rendyananta/example-online-book-store/internal/presenter/http/user"
	bookrp "github.com/rendyananta/example-online-book-store/internal/repo/book"
	orderrp "github.com/rendyananta/example-online-book-store/internal/repo/order"
	userrp "github.com/rendyananta/example-online-book-store/internal/repo/user"
	bookuc "github.com/rendyananta/example-online-book-store/internal/usecase/book"
	orderuc "github.com/rendyananta/example-online-book-store/internal/usecase/order"
	useruc "github.com/rendyananta/example-online-book-store/internal/usecase/user"
	"github.com/rendyananta/example-online-book-store/pkg/auth"
	"github.com/rendyananta/example-online-book-store/pkg/cache"
	"github.com/rendyananta/example-online-book-store/pkg/db"
)

type GlobalModules struct {
	DBConnManager *db.ConnManager
	CacheManager  *cache.Manager
	AuthManager   *auth.Manager
}

type RepoModules struct {
	BookRepo  *bookrp.Repo
	UserRepo  *userrp.Repo
	OrderRepo *orderrp.Repo
}

type UseCaseModules struct {
	UserAuthentication *useruc.AuthenticatorUseCase
	UserRegistration   *useruc.RegisterUseCase
	BookQueries        *bookuc.QueriesUseCase
	OrderPlacement     *orderuc.PlaceOrderUseCase
	OrderQueries       *orderuc.QueriesUseCase
}

type HTTPHandlers struct {
	Auth  user.Handler
	Book  book.Handler
	Order order.Handler
}
