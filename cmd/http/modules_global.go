package main

import (
	"github.com/rendyananta/example-online-book-store/pkg/validator"
	"log/slog"

	"github.com/rendyananta/example-online-book-store/pkg/auth"
	"github.com/rendyananta/example-online-book-store/pkg/cache"
	"github.com/rendyananta/example-online-book-store/pkg/db"
	"github.com/rendyananta/example-online-book-store/pkg/log"
)

func loadGlobalModules(cfg BinaryConfig) GlobalModules {
	log.SetUp(cfg.App.Global.Log)
	validator.SetUp()

	dbManager, err := db.NewConnectionManager(cfg.App.Global.DB)
	if err != nil {
		slog.Error("cannot initialize db connection manager", slog.String("err", err.Error()))
		panic(err)
	}

	cacheManager := cache.NewManager(cfg.App.Global.Cache)

	cmDBDriver, err := cache.NewDatabaseDriver(cfg.App.Global.CacheDBDriver, dbManager)
	if err != nil {
		slog.Error("cannot initialize db connection manager", slog.String("err", err.Error()))
	}

	cacheManager.Register(cache.DrvNameDatabase, cmDBDriver)

	authManager, err := auth.NewAuthManager(cfg.App.Global.Auth, cacheManager)
	if err != nil {
		slog.Error("cannot initialize auth manager", slog.String("err", err.Error()))
		panic(err)
	}

	return GlobalModules{
		DBConnManager: dbManager,
		CacheManager:  &cacheManager,
		AuthManager:   authManager,
	}
}
