package migrations

import "github.com/jmoiron/sqlx"

type CreateAuthorsTable struct {
	Conn *sqlx.DB
}

func (c CreateAuthorsTable) Up() error {
	query := `create table if not exists authors (
                       id uuid primary key,
                       name varchar(255) not null,
                       created_at timestamp not null default current_timestamp,
                       updated_at timestamp not null default current_timestamp,
                       deleted_at timestamp
        )`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateAuthorsTable) Down() error {
	query := `drop table if exists authors`

	_, err := c.Conn.Exec(query)
	return err
}
