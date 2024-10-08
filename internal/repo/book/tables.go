package book

import (
	"database/sql"
)

type tableBook struct {
	ID               string       `db:"id"`
	Title            string       `db:"title"`
	Description      string       `db:"description"`
	Price            float64      `db:"price"`
	ISBN             string       `db:"isbn"`
	PublisherID      string       `db:"publisher_id"`
	PublisherName    string       `db:"publisher_name"`
	Language         string       `db:"language"`
	Edition          string       `db:"edition"`
	Pages            int          `db:"pages"`
	PublishedAt      sql.NullTime `db:"published_at"`
	FirstPublishedAt sql.NullTime `db:"first_published_at"`
	CoverImg         string       `db:"cover_img"`
	Rating           float64      `db:"rating"`
	CreatedAt        sql.NullTime `db:"created_at"`
	UpdatedAt        sql.NullTime `db:"updated_at"`
}

type tableAuthor struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	BookID string `db:"book_id"`
}

type tableGenre struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	BookID string `db:"book_id"`
}
