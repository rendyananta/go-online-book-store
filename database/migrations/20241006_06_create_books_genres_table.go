package migrations

import "github.com/jmoiron/sqlx"

type CreateBooksGenresTable struct {
	Conn *sqlx.DB
}

func (c CreateBooksGenresTable) Up() error {
	query := `create table if not exists books_genres (
                       id uuid primary key,
                       book_id uuid not null,
                       genre_id uuid not null
        );
		create index if not exists books_genres_book_id_idx on books_genres (book_id);
		create index if not exists books_genres_genre_id_idx on books_genres (genre_id);
        `

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateBooksGenresTable) Down() error {
	query := `drop table if exists books_genres;
			drop index if exists books_genres_book_id_idx;
			drop index if exists books_genres_genre_id_idx;
	`

	_, err := c.Conn.Exec(query)
	return err
}
