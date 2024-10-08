package main

import (
	"flag"
	"fmt"
	"github.com/rendyananta/example-online-book-store/database/migrations"
	"github.com/rendyananta/example-online-book-store/internal/config"
	"github.com/rendyananta/example-online-book-store/pkg/db"
	"github.com/rendyananta/example-online-book-store/pkg/log"
	"log/slog"
	"reflect"
)

type Migration interface {
	Up() error
	Down() error
}

// simple migrations binary
func main() {
	cfg := BinaryConfig{
		App: config.LoadAppConfig(),
	}

	var upCmd = true

	log.SetUp(cfg.App.Global.Log)

	flag.Parse()
	args := flag.Args()

	slog.Info(fmt.Sprintf("args: %s", args))

	if len(args) > 0 {
		switch args[0] {
		case "up":
			upCmd = true
		case "down":
			upCmd = false
		}
	}

	dbManager, err := db.NewConnectionManager(cfg.App.Global.DB)
	if err != nil {
		slog.Error("cannot initialize db connection manager", slog.String("err", err.Error()))
		panic(err)
	}

	defaultConn, err := dbManager.Connection(db.ConnDefault)
	if err != nil {
		slog.Error("cannot get connection")
	}

	pendingMigrations := []Migration{
		&migrations.CreateUsersTable{Conn: defaultConn},
		&migrations.CreatePublishersTable{Conn: defaultConn},
		&migrations.CreateAuthorsTable{Conn: defaultConn},
		&migrations.CreateGenresTable{Conn: defaultConn},
		&migrations.CreateBooksTable{Conn: defaultConn},
		&migrations.CreateBooksAuthorsTable{Conn: defaultConn},
		&migrations.CreateBooksGenresTable{Conn: defaultConn},
		&migrations.CreateOrdersTable{Conn: defaultConn},
		&migrations.CreateOrderLinesTable{Conn: defaultConn},
	}

	if upCmd {
		for _, migration := range pendingMigrations {
			t := reflect.TypeOf(migration)
			migrationName := t.Elem().Name()

			slog.Info(fmt.Sprintf("migrating [%s] | UP", migrationName))

			err = migration.Up()
			if err != nil {
				slog.Error("cannot run migration", slog.String("err", err.Error()))
				break
			}
		}

		slog.Info("all up migrations done")

		return
	}

	for i := len(pendingMigrations) - 1; i >= 0; i-- {
		migration := pendingMigrations[i]
		t := reflect.TypeOf(migration)
		migrationName := t.Elem().Name()

		slog.Info(fmt.Sprintf("migrating [%s] | DOWN", migrationName))

		err = migration.Down()
		if err != nil {
			slog.Error("cannot run migration", slog.String("err", err.Error()))
			break
		}
	}

	slog.Info("all down migrations done")
}
