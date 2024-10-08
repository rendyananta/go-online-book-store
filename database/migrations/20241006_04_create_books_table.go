package migrations

import "github.com/jmoiron/sqlx"

type CreateBooksTable struct {
	Conn *sqlx.DB
}

func (c CreateBooksTable) Up() error {
	query := `create table if not exists books (
                       id uuid primary key,
                       title varchar (255) not null,
                       description text not null,
                       price double not null,
                       isbn varchar not null,
                       language varchar(100),
                       edition varchar(255),
                       pages int,
                       publisher_id uuid not null,
                       published_at timestamp,
                       first_published_at timestamp,
                       cover_img text,
                       rating double,
                       created_at timestamp not null default current_timestamp,
                       updated_at timestamp not null default current_timestamp,
                       deleted_at timestamp
        );

-- 		create virtual table if not exists book_search_index using fts5 (
-- 		    id, title, description, genres, authors, publisher
-- 		);
`

	_, err := c.Conn.Exec(query)
	return err
}

func (c CreateBooksTable) Down() error {
	query := `drop table if exists books;
--  		drop table book_search_index;
		`

	_, err := c.Conn.Exec(query)
	return err
}
