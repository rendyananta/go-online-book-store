package main

import (
	"log/slog"

	bookuc "github.com/rendyananta/example-online-book-store/internal/usecase/book"
	orderuc "github.com/rendyananta/example-online-book-store/internal/usecase/order"
	useruc "github.com/rendyananta/example-online-book-store/internal/usecase/user"
)

func loadUseCaseModules(_ BinaryConfig, globalModules GlobalModules, repoModules RepoModules) UseCaseModules {
	userAuthentication, err := useruc.NewAuthenticatorUseCase(repoModules.UserRepo, globalModules.AuthManager)
	if err != nil {
		slog.Error("cannot initialize user auth use case", slog.String("err", err.Error()))
		panic(err)
	}

	userRegistration, err := useruc.NewRegisterUseCase(repoModules.UserRepo)
	if err != nil {
		slog.Error("cannot initialize user registration use case", slog.String("err", err.Error()))
		panic(err)
	}

	bookQueries, err := bookuc.NewQueryUseCase(repoModules.BookRepo)
	if err != nil {
		slog.Error("cannot initialize book queries use case", slog.String("err", err.Error()))
		panic(err)
	}

	orderPlacement, err := orderuc.NewPlaceOrderUseCase(repoModules.OrderRepo, repoModules.BookRepo)
	if err != nil {
		slog.Error("cannot initialize place order use case", slog.String("err", err.Error()))
		panic(err)
	}

	orderQueries, err := orderuc.NewOrderQueriesUseCase(repoModules.OrderRepo, repoModules.BookRepo)
	if err != nil {
		slog.Error("cannot initialize place queries use case", slog.String("err", err.Error()))
		panic(err)
	}

	return UseCaseModules{
		UserAuthentication: userAuthentication,
		UserRegistration:   userRegistration,
		BookQueries:        bookQueries,
		OrderPlacement:     orderPlacement,
		OrderQueries:       orderQueries,
	}
}
