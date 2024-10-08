package order

import (
	"database/sql"
)

type tableOrder struct {
	ID         string       `db:"id"`
	UserID     string       `db:"user_id"`
	GrandTotal float64      `db:"grand_total"`
	Status     string       `db:"status"`
	CreatedAt  sql.NullTime `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
}

type tableOrderLine struct {
	ID                string  `db:"id"`
	OrderID           string  `db:"order_id"`
	LineReferenceType string  `db:"line_reference_type"`
	LineReferenceID   string  `db:"line_reference_id"`
	Amount            float64 `db:"amount"`
	Quantity          int     `db:"quantity"`
	Subtotal          float64 `db:"subtotal"`
}
