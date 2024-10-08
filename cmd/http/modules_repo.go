package main

import (
	"log/slog"

	"github.com/rendyananta/example-online-book-store/internal/repo/book"
	"github.com/rendyananta/example-online-book-store/internal/repo/order"
	"github.com/rendyananta/example-online-book-store/internal/repo/user"
)

func loadRepoModules(cfg BinaryConfig, globalModules GlobalModules) RepoModules {
	bookRepo, err := book.NewBookRepo(cfg.App.Domain.BookRepo, globalModules.DBConnManager)
	if err != nil {
		slog.Error("cannot initialize book repo", slog.String("err", err.Error()))
		panic(err)
	}

	userRepo, err := user.NewUserRepo(cfg.App.Domain.UserRepo, globalModules.DBConnManager)
	if err != nil {
		slog.Error("cannot initialize book repo", slog.String("err", err.Error()))
		panic(err)
	}

	orderRepo, err := order.NewOrderRepo(cfg.App.Domain.OrderRepo, globalModules.DBConnManager)
	if err != nil {
		slog.Error("cannot initialize book repo", slog.String("err", err.Error()))
		panic(err)
	}

	return RepoModules{
		BookRepo:  bookRepo,
		UserRepo:  userRepo,
		OrderRepo: orderRepo,
	}
}
