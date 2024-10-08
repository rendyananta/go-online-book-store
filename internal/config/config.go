package config

import (
	bookrp "github.com/rendyananta/example-online-book-store/internal/repo/book"
	orderrp "github.com/rendyananta/example-online-book-store/internal/repo/order"
	userrp "github.com/rendyananta/example-online-book-store/internal/repo/user"
	"github.com/rendyananta/example-online-book-store/pkg/auth"
	"github.com/rendyananta/example-online-book-store/pkg/cache"
	"github.com/rendyananta/example-online-book-store/pkg/db"
	"github.com/rendyananta/example-online-book-store/pkg/log"
)

type App struct {
	Global Global
	Domain Domain
}

type Global struct {
	Log           log.Config
	DB            db.Config
	Cache         cache.Config
	CacheDBDriver cache.DriverDatabaseConfig
	Auth          auth.Config
}

type Domain struct {
	UserRepo  userrp.Config
	BookRepo  bookrp.Config
	OrderRepo orderrp.Config
}
