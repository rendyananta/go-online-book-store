package main

import (
	"github.com/rendyananta/example-online-book-store/internal/presenter/http"
	"github.com/rendyananta/example-online-book-store/internal/presenter/http/book"
	"github.com/rendyananta/example-online-book-store/internal/presenter/http/order"
	"github.com/rendyananta/example-online-book-store/internal/presenter/http/user"
	"github.com/rendyananta/example-online-book-store/pkg/auth"
)

func loadHTTPHandlers(_ BinaryConfig, globalModules GlobalModules, repoModules RepoModules, useCaseModules UseCaseModules) HTTPHandlers {
	authMiddleware := auth.NewMiddleware(globalModules.AuthManager, &http.AppResponseWriter{})

	return HTTPHandlers{
		Auth: user.Handler{
			Register:      useCaseModules.UserRegistration,
			Authenticator: useCaseModules.UserAuthentication,
		},
		Book: book.Handler{
			Queries: useCaseModules.BookQueries,
		},
		Order: order.Handler{
			AuthMiddleware:    authMiddleware,
			PlaceOrderUseCase: useCaseModules.OrderPlacement,
			Queries:           useCaseModules.OrderQueries,
		},
	}
}
