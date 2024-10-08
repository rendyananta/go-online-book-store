package book

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"github.com/rendyananta/example-online-book-store/pkg/db"
	"log/slog"
	"sync"
	"time"
)

//go:generate mockgen -source=book.go -destination=book_db_conn_mock_test.go -package book
type dbConnManager interface {
	Connection(name string) (*sqlx.DB, error)
}

type dbConnection interface {
	Preparex(query string) (*sqlx.Stmt, error)
	Rebind(query string) string

	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type preparedQueryGetter interface {
	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error
}

type Config struct {
	DBConn string
}

type preparedStmt struct {
	bookPagination       preparedQueryGetter
	bookSearchPagination preparedQueryGetter
}

type Repo struct {
	cfg          Config
	dbConn       dbConnection
	preparedStmt preparedStmt
}

func NewBookRepo(cfg Config, dbConnManager dbConnManager) (*Repo, error) {
	if cfg.DBConn == "" {
		cfg.DBConn = db.ConnDefault
	}

	conn, err := dbConnManager.Connection(cfg.DBConn)
	if err != nil {
		return nil, err
	}

	r := &Repo{
		dbConn: conn,
	}

	if err = r.boot(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Repo) boot() error {
	var err error

	r.preparedStmt.bookPagination, err = r.dbConn.Preparex(r.dbConn.Rebind(queryPaginateAllBooks))
	if err != nil {
		return err
	}

	r.preparedStmt.bookSearchPagination, err = r.dbConn.Preparex(r.dbConn.Rebind(queryPaginateBooksSearchResult))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) compileBooksResult(ctx context.Context, result []tableBook) ([]book.Book, error) {
	var books = make([]book.Book, 0, len(result))
	var bookIdxByID = make(map[string]int)

	var bookIDs = make([]string, 0, len(result))
	for i, item := range result {
		bookIDs = append(bookIDs, item.ID)

		var publishedAt *time.Time
		var firstPublishedAt *time.Time

		if item.PublishedAt.Valid {
			publishedAt = &item.PublishedAt.Time
		}

		if item.FirstPublishedAt.Valid {
			firstPublishedAt = &item.FirstPublishedAt.Time
		}

		books = append(books, book.Book{
			ID:               item.ID,
			Title:            item.Title,
			Description:      item.Description,
			Price:            item.Price,
			ISBN:             item.ISBN,
			Language:         item.Language,
			Edition:          item.Edition,
			Pages:            item.Pages,
			PublishedAt:      publishedAt,
			FirstPublishedAt: firstPublishedAt,
			CoverImg:         item.CoverImg,
			Rating:           item.Rating,
			Publisher: book.Publisher{
				ID:   item.PublisherID,
				Name: item.PublisherName,
			},
			Authors: nil,
			Genres:  nil,
		})

		bookIdxByID[item.ID] = i
	}

	var wg sync.WaitGroup
	var errs []error

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		authors, err := r.getAuthorsByBookIDs(ctx, bookIDs)
		if err != nil {
			errs = append(errs, err)
			return
		}

		for _, author := range authors {
			bookIdx, ok := bookIdxByID[author.BookID]
			if !ok {
				continue
			}

			books[bookIdx].Authors = append(books[bookIdx].Authors, book.Author{
				ID:   author.ID,
				Name: author.Name,
			})
		}
	}(ctx, &wg)

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		genres, err := r.getGenresByBookIDs(ctx, bookIDs)
		if err != nil {
			errs = append(errs, err)
			return
		}

		for _, genre := range genres {
			bookIdx, ok := bookIdxByID[genre.BookID]
			if !ok {
				continue
			}

			books[bookIdx].Genres = append(books[bookIdx].Genres, book.Genre{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}
	}(ctx, &wg)
	wg.Wait()

	if len(errs) > 0 {
		return books, errs[0]
	}

	return books, nil
}

func (r *Repo) PaginateBookSearch(ctx context.Context, searchQuery string, param book.PaginationParam) (book.PaginationResult, error) {
	if param.PerPage == 0 {
		param.PerPage = defaultPaginationLength
	}

	var result []tableBook

	searchQuery = fmt.Sprintf("%s%s%s", "%", searchQuery, "%")

	err := r.preparedStmt.bookSearchPagination.SelectContext(ctx, &result, searchQuery, param.LastID, searchQuery, param.PerPage)
	if err != nil {
		return book.PaginationResult{}, err
	}

	books, err := r.compileBooksResult(ctx, result)
	lastID := ""
	if len(books) > 0 {
		lastID = books[len(books)-1].ID
	}

	if err != nil {
		return book.PaginationResult{
			Data:    books,
			PerPage: param.PerPage,
			LastID:  lastID,
		}, err
	}

	return book.PaginationResult{
		Data:    books,
		PerPage: param.PerPage,
		LastID:  lastID,
	}, nil
}

func (r *Repo) PaginateAllBooks(ctx context.Context, param book.PaginationParam) (book.PaginationResult, error) {
	if param.PerPage == 0 {
		param.PerPage = defaultPaginationLength
	}

	var result []tableBook

	err := r.preparedStmt.bookPagination.SelectContext(ctx, &result, param.LastID, param.PerPage)
	if err != nil {
		return book.PaginationResult{}, err
	}

	books, err := r.compileBooksResult(ctx, result)

	lastID := ""
	if len(books) > 0 {
		lastID = books[len(books)-1].ID
	}

	if err != nil {
		return book.PaginationResult{
			Data:    books,
			PerPage: param.PerPage,
			LastID:  lastID,
		}, err
	}

	return book.PaginationResult{
		Data:    books,
		PerPage: param.PerPage,
		LastID:  lastID,
	}, nil
}

func (r *Repo) getAuthorsByBookIDs(ctx context.Context, bookIDs []string) ([]tableAuthor, error) {
	var authors []tableAuthor

	query, args, err := sqlx.In(queryGetAuthorsByBookIDs, bookIDs)
	if err != nil {
		slog.Error("error preparing authors in query", slog.String("error", err.Error()))
		return nil, err
	}

	err = r.dbConn.SelectContext(ctx, &authors, r.dbConn.Rebind(query), args...)
	if err != nil {
		slog.Error("error executing authors by book ids query", slog.String("error", err.Error()))
		return nil, err
	}

	return authors, nil
}

func (r *Repo) getGenresByBookIDs(ctx context.Context, bookIDs []string) ([]tableGenre, error) {
	var genres []tableGenre

	query, args, err := sqlx.In(queryGetGenresByBookIDs, bookIDs)
	if err != nil {
		slog.Error("error preparing genres in query", slog.String("error", err.Error()))
		return nil, err
	}

	err = r.dbConn.SelectContext(ctx, &genres, r.dbConn.Rebind(query), args...)
	if err != nil {
		slog.Error("error executing genres by book ids query", slog.String("error", err.Error()))
		return nil, err
	}

	return genres, nil
}

func (r *Repo) FindByIDs(ctx context.Context, bookIDs []string) ([]book.Book, error) {
	var rawBooks []tableBook

	query, args, err := sqlx.In(queryGetBookByIDs, bookIDs)
	if err != nil {
		return nil, err
	}

	err = r.dbConn.SelectContext(ctx, &rawBooks, r.dbConn.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	var books = make([]book.Book, 0, len(rawBooks))
	var bookIdxByID = make(map[string]int)
	var resultBookIDs = make([]string, 0, len(rawBooks))

	for i, itemResult := range rawBooks {
		var publishedAt *time.Time
		var firstPublishedAt *time.Time

		if itemResult.PublishedAt.Valid {
			publishedAt = &itemResult.PublishedAt.Time
		}

		if itemResult.FirstPublishedAt.Valid {
			firstPublishedAt = &itemResult.FirstPublishedAt.Time
		}

		books = append(books, book.Book{
			ID:               itemResult.ID,
			Title:            itemResult.Title,
			Description:      itemResult.Description,
			Price:            itemResult.Price,
			ISBN:             itemResult.ISBN,
			Language:         itemResult.Language,
			Edition:          itemResult.Edition,
			Pages:            itemResult.Pages,
			PublishedAt:      publishedAt,
			FirstPublishedAt: firstPublishedAt,
			CoverImg:         itemResult.CoverImg,
			Rating:           itemResult.Rating,
			Publisher: book.Publisher{
				ID:   itemResult.PublisherID,
				Name: itemResult.PublisherName,
			},
			Authors: nil,
			Genres:  nil,
		})

		bookIdxByID[itemResult.ID] = i
		resultBookIDs = append(resultBookIDs, itemResult.ID)
	}

	var wg sync.WaitGroup
	var errs []error

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		authors, err := r.getAuthorsByBookIDs(ctx, resultBookIDs)
		if err != nil {
			errs = append(errs, err)
			return
		}

		if len(authors) == 0 {
			return
		}

		for _, author := range authors {
			bookIdx, ok := bookIdxByID[author.BookID]
			if !ok {
				continue
			}

			books[bookIdx].Authors = append(books[bookIdx].Authors, book.Author{
				ID:   author.ID,
				Name: author.Name,
			})
		}
	}(ctx, &wg)

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		genres, err := r.getGenresByBookIDs(ctx, resultBookIDs)
		if err != nil {
			errs = append(errs, err)
			return
		}

		if len(genres) == 0 {
			return
		}

		for _, genre := range genres {
			bookIdx, ok := bookIdxByID[genre.BookID]
			if !ok {
				continue
			}

			books[bookIdx].Genres = append(books[bookIdx].Genres, book.Genre{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}
	}(ctx, &wg)
	wg.Wait()

	if len(errs) > 0 {
		return books, errs[0]
	}

	return books, nil
}
