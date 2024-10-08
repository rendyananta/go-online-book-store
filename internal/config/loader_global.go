package config

import (
	"github.com/rendyananta/example-online-book-store/pkg/auth"
	"github.com/rendyananta/example-online-book-store/pkg/cache"
	"github.com/rendyananta/example-online-book-store/pkg/db"
	"github.com/rendyananta/example-online-book-store/pkg/log"
)

func loadGlobalConfig() Global {
	return Global{
		Log: log.Config{
			LogPath:       LoadFromEnvString("LOG_PATH", ""),
			JSONFormatted: LoadFromEnvBool("LOG_JSON_FORMATTED", false),
		},
		DB: db.Config{
			Connections: map[string]db.ConnectionConfig{
				db.ConnDefault: {
					DSN:        LoadFromEnvString("DB_DEFAULT_DSN", "file:database/app.sqlite3"),
					DriverName: LoadFromEnvString("DB_DEFAULT_DRIVER", "sqlite3"),
				},
				db.ConnCache: {
					DSN:        LoadFromEnvString("DB_CACHE_DSN", "file:database/cache.sqlite3"),
					DriverName: LoadFromEnvString("DB_CACHE_DRIVER", "sqlite3"),
				},
			},
		},
		Cache: cache.Config{
			DefaultDriver: LoadFromEnvString("CACHE_DRIVER", cache.DrvNameDatabase),
		},
		Auth: auth.Config{
			TokenLifetime: LoadFromEnvTimeDuration("AUTH_TOKEN_LIFETIME", 0),
			CipherKeys:    LoadFromEnvStringSlice("AUTH_CIPHER_KEYS", nil),
		},
	}
}
