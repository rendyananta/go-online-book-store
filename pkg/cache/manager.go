package cache

import (
	"context"
	"time"
)

type Driver interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type Config struct {
	DefaultDriver string
}

type DriverName = string

type Manager struct {
	driverSet map[DriverName]Driver
	config    Config
}

func (m Manager) Register(name DriverName, driver Driver) {
	m.driverSet[name] = driver
}

func NewManager(cfg Config) Manager {
	if cfg.DefaultDriver == "" {
		cfg.DefaultDriver = DrvNameDatabase
	}

	return Manager{
		config:    cfg,
		driverSet: make(map[DriverName]Driver),
	}
}

func (m Manager) Get(ctx context.Context, key string) ([]byte, error) {
	drv, registered := m.driverSet[m.config.DefaultDriver]
	if !registered {
		return nil, ErrDriverUnregistered
	}

	return drv.Get(ctx, key)
}

func (m Manager) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	drv, registered := m.driverSet[m.config.DefaultDriver]
	if !registered {
		return ErrDriverUnregistered
	}

	return drv.Set(ctx, key, val, ttl)
}

func (m Manager) Del(ctx context.Context, key string) error {
	drv, registered := m.driverSet[m.config.DefaultDriver]
	if !registered {
		return ErrDriverUnregistered
	}

	return drv.Del(ctx, key)
}
