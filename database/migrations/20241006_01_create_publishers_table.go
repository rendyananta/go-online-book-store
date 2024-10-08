package migrations

import "github.com/jmoiron/sqlx"

type CreatePublishersTable struct {
	Conn *sqlx.DB
}

func (c CreatePublishersTable) Up() error {
	query := `create table if not exists publishers (
                       id uuid primary key,
                       name varchar(255) not null,
                       created_at timestamp not null default current_timestamp,
                       updated_at timestamp not null default current_timestamp,
                       deleted_at timestamp
        )`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreatePublishersTable) Down() error {
	query := `drop table if exists publishers`

	_, err := c.Conn.Exec(query)
	return err
}
