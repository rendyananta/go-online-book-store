package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	"github.com/rendyananta/example-online-book-store/pkg/db"
)

//go:generate mockgen -source=user.go -destination=user_db_conn_mock_test.go -package user
type dbConnManager interface {
	Connection(name string) (*sqlx.DB, error)
}

type dbConnection interface {
	Preparex(query string) (*sqlx.Stmt, error)
	Rebind(query string) string
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type queryGetter interface {
	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
}

type queryExecer interface {
	ExecContext(ctx context.Context, args ...any) (sql.Result, error)
}

type Config struct {
	DBConn string
}

type preparedStmt struct {
	findByEmailStmt queryGetter
	findByIDStmt    queryGetter
}

type Repo struct {
	cfg          Config
	dbConn       dbConnection
	preparedStmt preparedStmt
}

func NewUserRepo(cfg Config, dbConnManager dbConnManager) (*Repo, error) {
	if cfg.DBConn == "" {
		cfg.DBConn = db.ConnDefault
	}

	conn, err := dbConnManager.Connection(cfg.DBConn)
	if err != nil {
		return nil, err
	}

	r := &Repo{
		dbConn: conn,
	}

	if err = r.boot(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Repo) boot() error {
	var err error

	r.preparedStmt.findByEmailStmt, err = r.dbConn.Preparex(r.dbConn.Rebind(queryGetUserByEmail))
	if err != nil {
		return err
	}

	r.preparedStmt.findByIDStmt, err = r.dbConn.Preparex(r.dbConn.Rebind(queryGetUserByID))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindByEmail(ctx context.Context, email string) (user.User, error) {
	var userResult tableUser
	err := r.preparedStmt.findByEmailStmt.GetContext(ctx, &userResult, email)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user.User{}, ErrEmailIsNotRegistered
	}

	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:       userResult.ID,
		Name:     userResult.Name,
		Email:    userResult.Email,
		Password: userResult.Password,
	}, nil
}

func (r *Repo) FindByID(ctx context.Context, id string) (user.User, error) {
	var userResult tableUser
	err := r.preparedStmt.findByEmailStmt.GetContext(ctx, &userResult, id)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user.User{}, ErrNotFound
	}

	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:       userResult.ID,
		Name:     userResult.Name,
		Email:    userResult.Email,
		Password: userResult.Password,
	}, nil
}

func (r *Repo) Create(ctx context.Context, param user.User) (user.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return user.User{}, err
	}

	now := time.Now()

	_, err = r.dbConn.ExecContext(ctx, queryInsertUser, id.String(), param.Name, param.Email, param.Password, now, now)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:       id.String(),
		Name:     param.Name,
		Email:    param.Email,
		Password: param.Password,
	}, nil
}
