package migrations

import "github.com/jmoiron/sqlx"

type CreateBooksAuthorsTable struct {
	Conn *sqlx.DB
}

func (c CreateBooksAuthorsTable) Up() error {
	query := `create table if not exists books_authors (
                       id uuid primary key,
                       book_id uuid not null,
                       author_id uuid not null
        );
		create index if not exists books_authors_book_id_idx on books_authors (book_id);
		create index if not exists books_authors_genre_id_idx on books_authors (author_id);
		`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateBooksAuthorsTable) Down() error {
	query := `drop table if exists books_authors;
			drop index if exists books_genres_book_id_idx;
			drop index if exists books_genres_genre_id_idx;
`

	_, err := c.Conn.Exec(query)
	return err
}
