package order

import (
	"time"

	"github.com/rendyananta/example-online-book-store/internal/entity/book"
)

type PaginationParam struct {
	PerPage int
	LastID  string
}

type PaginationResult struct {
	Data    []Main
	PerPage int
	LastID  string
}

type Status = string

const (
	StatusCreated Status = "created"
	StatusDone    Status = "done"
)

type LineReferenceType = string
type LineReference interface {
	book.Book
}

const (
	LineReferenceTypeBook LineReferenceType = "book"
)

type Main struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	GrandTotal float64    `json:"grand_total"`
	Status     Status     `json:"status"`
	Lines      []Line     `json:"lines"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type Line struct {
	ID                string  `json:"id"`
	OrderID           string  `json:"order_id"`
	LineReferenceType string  `json:"line_reference_type"`
	LineReferenceID   string  `json:"line_reference_id"`
	LineItem          any     `json:"line_item"`
	Amount            float64 `json:"amount"`
	Quantity          int     `json:"quantity"`
	Subtotal          float64 `json:"subtotal"`
}
