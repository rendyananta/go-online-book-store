package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := loadConfig()
	globalModules := loadGlobalModules(cfg)

	slog.Info("starting app-http")
	slog.Info("global modules loaded")

	repoModules := loadRepoModules(cfg, globalModules)

	slog.Info("repo modules loaded")

	usecaseModules := loadUseCaseModules(cfg, globalModules, repoModules)

	slog.Info("use cases modules loaded")

	handlers := loadHTTPHandlers(cfg, globalModules, repoModules, usecaseModules)

	slog.Info("http handlers loaded")

	mux := http.NewServeMux()

	handlers.Auth.Handle(mux)
	handlers.Book.Handle(mux)
	handlers.Order.Handle(mux)

	slog.Info(fmt.Sprintf("listening http server on :%d", cfg.HTTP.ListenPort))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.ListenPort),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
			return
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	<-signalChan

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
