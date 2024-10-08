package migrations

import "github.com/jmoiron/sqlx"

type CreateOrdersTable struct {
	Conn *sqlx.DB
}

func (c CreateOrdersTable) Up() error {
	query := `create table if not exists orders (
                       id uuid primary key,
                       user_id uuid not null,
                       grand_total double not null,
                       status varchar(255) not null,
                       created_at timestamp not null default current_timestamp,
                       updated_at timestamp not null default current_timestamp,
                       deleted_at timestamp
        )`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateOrdersTable) Down() error {
	query := `drop table if exists orders`

	_, err := c.Conn.Exec(query)
	return err
}
