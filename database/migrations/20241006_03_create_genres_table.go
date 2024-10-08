package migrations

import "github.com/jmoiron/sqlx"

type CreateGenresTable struct {
	Conn *sqlx.DB
}

func (c CreateGenresTable) Up() error {
	query := `create table if not exists genres (
                       id uuid primary key,
                       name varchar(255) not null,
                       created_at timestamp not null default current_timestamp,
                       updated_at timestamp not null default current_timestamp
        )`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateGenresTable) Down() error {
	query := `drop table if exists genres`

	_, err := c.Conn.Exec(query)
	return err
}
