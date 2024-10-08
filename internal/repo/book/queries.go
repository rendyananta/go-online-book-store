package book

const (
	queryGetBookByIDs = `select b.id, title, description, price, isbn, language,edition, pages, publisher_id,  p.name as publisher_name,  published_at, first_published_at, cover_img, rating, b.created_at, b.updated_at 
								from books b
								inner join publishers p on p.id = b.publisher_id
								where b.id in (?) and b.deleted_at is null`

	queryGetAuthorsByBookIDs = `select a.id, a.name, b.id as book_id from books_authors 
    							inner join books b on b.id = books_authors.book_id 
    							inner join authors a on books_authors.author_id = a.id
    							where books_authors.book_id in (?)`

	queryGetGenresByBookIDs = `select g.id, g.name, b.id as book_id from books_genres
								inner join books b on b.id = books_genres.book_id
								inner join genres g on books_genres.genre_id = g.id
    							where books_genres.book_id in (?)`

	queryPaginateAllBooks = `select b.id, title, description, price, isbn, language,edition, pages, publisher_id, p.name as publisher_name, 
       							 published_at, first_published_at, cover_img, rating, b.created_at, b.updated_at 
									from books b
									inner join publishers p on p.id = b.publisher_id
									where b.deleted_at is null and b.id > ?
									order by b.id limit ?`

	//queryPaginateBooksSearchResult = `select b.id, title, description, price, isbn, language,edition, pages, publisher_id, p.name as publisher_name,
	//   							 published_at, first_published_at, cover_img, rating, b.created_at, b.updated_at
	//								from book_search_index bsi
	//								inner join books b on b.id = bsi.id
	//								inner join publishers p on p.id = b.publisher_id
	//								where book_search_index match ?
	//								  and b.deleted_at is null
	//								  and b.id > ?
	//								order by rank(), b.id limit ?`

	queryPaginateBooksSearchResult = `select b.id, title, description, price, isbn, language,edition, pages, publisher_id, p.name as publisher_name, 
       							 published_at, first_published_at, cover_img, rating, b.created_at, b.updated_at 
									from books b
									inner join publishers p on p.id = b.publisher_id
									where b.title like ?
										and b.deleted_at is null and b.id > ?
									order by b.title like ?, b.id limit ?`
)
