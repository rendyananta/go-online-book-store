package migrations

import "github.com/jmoiron/sqlx"

type CreateOrderLinesTable struct {
	Conn *sqlx.DB
}

func (c CreateOrderLinesTable) Up() error {
	query := `create table if not exists order_lines (
                       id uuid primary key,
                       order_id uuid not null,
                       line_reference_type varchar(255) not null,
                       line_reference_id uuid not null, -- book / shipping fee / discount / platform fee
                       amount double not null,
                       quantity int not null default 1,
                       subtotal double not null 
        )`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateOrderLinesTable) Down() error {
	query := `drop table if exists order_lines`

	_, err := c.Conn.Exec(query)
	return err
}
