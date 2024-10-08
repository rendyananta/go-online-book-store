package migrations

import (
	"github.com/jmoiron/sqlx"
)

type CreateUsersTable struct {
	Conn *sqlx.DB
}

func (c CreateUsersTable) Up() error {
	query := `create table if not exists users (
                       id uuid primary key,
                       name varchar (255) not null,
                       email varchar(255) not null unique,
                       password varchar(255) not null,
                       created_at timestamp not null default current_timestamp,
                       updated_at timestamp not null default current_timestamp
        )`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateUsersTable) Down() error {
	query := `drop table if exists users`

	_, err := c.Conn.Exec(query)
	return err
}
