package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ConnectionConfig struct {
	DSN        string
	DriverName string
}

type Config struct {
	Connections map[string]ConnectionConfig
}

type ConnManager struct {
	connections map[string]*sqlx.DB
}

func NewConnectionManager(cfg Config) (*ConnManager, error) {
	var connections = make(map[string]*sqlx.DB)

	if len(cfg.Connections) == 0 {
		cfg.Connections[ConnDefault] = ConnectionConfig{
			DSN:        ":memory:",
			DriverName: "sqlite3",
		}
	}

	for s, config := range cfg.Connections {
		conn, err := sqlx.Open(config.DriverName, config.DSN)
		if err != nil {
			return nil, err
		}

		connections[s] = conn
	}

	return &ConnManager{
		connections: connections,
	}, nil
}

func (m ConnManager) Connection(name string) (*sqlx.DB, error) {
	conn, exist := m.connections[name]
	if !exist {
		return nil, ErrConnectionUnregistered
	}

	return conn, nil
}
