package cache

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type dbConnManager interface {
	Connection(name string) (*sqlx.DB, error)
}

const DrvNameDatabase DriverName = "database"
const drvDatabaseDefaultConn = "default"
const drvDatabaseDefaultTrimDuration = 5 * time.Minute

const (
	querySetKey = `insert into caches (id, key, value, expired_at) 
		values (?, ?, ?, ?) on conflict (key) do update set value = ?, expired_at = ?`
	queryDelKey = `delete from caches where key = ?`
)

type DriverDatabaseConfig struct {
	TrimDuration time.Duration
	Connection   string
}

type DriverDatabase struct {
	connection      *sqlx.DB
	config          DriverDatabaseConfig
	getPreparedStmt *sqlx.Stmt
}

func NewDatabaseDriver(config DriverDatabaseConfig, connManager dbConnManager) (*DriverDatabase, error) {
	if config.Connection == "" {
		config.Connection = drvDatabaseDefaultConn
	}

	if config.TrimDuration <= 0 {
		config.TrimDuration = drvDatabaseDefaultTrimDuration
	}

	conn, err := connManager.Connection(config.Connection)
	if err != nil {
		return nil, err
	}

	drv := &DriverDatabase{
		connection: conn,
		config:     config,
	}

	if err := drv.boot(); err != nil {
		return nil, err
	}

	return drv, nil
}

func (d *DriverDatabase) boot() error {
	_, err := d.connection.Exec(d.connection.Rebind(`create table if not exists caches (id varchar(36) primary key, key string unique, value text, expired_at timestamp);
		create index if not exists key_expired_at_idx on caches (key, expired_at);
		create index if not exists expired_at on caches (expired_at);
		`))
	if err != nil {
		return err
	}

	d.getPreparedStmt, err = d.connection.Preparex(d.connection.Rebind("select value from caches where key = ? and expired_at > ?"))
	if err != nil {
		return err
	}

	ticker := time.NewTicker(d.config.TrimDuration)

	go func() {
		defer ticker.Stop()

		for {
			<-ticker.C
			if _, cleanupErr := d.connection.Query("delete from caches where expired_at > ?", time.Now()); cleanupErr != nil {
				slog.Error(fmt.Sprintf("failed to trim caches table, err: %s", cleanupErr))
			}
		}
	}()

	return nil
}

func (d *DriverDatabase) Get(ctx context.Context, key string) ([]byte, error) {
	var result []byte

	err := d.getPreparedStmt.GetContext(ctx, &result, key, time.Now())

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return result, nil
}

func (d *DriverDatabase) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	var expiredAt *time.Time

	if ttl > 0 {
		expiryTime := time.Now().Add(ttl)
		expiredAt = &expiryTime
	}

	newVal := bytes.Clone(val)

	query := d.connection.Rebind(querySetKey)

	if _, err = d.connection.ExecContext(ctx, query, id.String(), key, val, expiredAt, newVal, expiredAt); err != nil {
		return err
	}

	return nil
}

func (d *DriverDatabase) Del(ctx context.Context, key string) error {
	if _, err := d.connection.ExecContext(ctx, d.connection.Rebind(queryDelKey), key); err != nil {
		return err
	}

	return nil
}
