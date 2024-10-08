package book

import (
	"context"
	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	bookrp "github.com/rendyananta/example-online-book-store/internal/repo/book"
)

//go:generate mockgen -source=query.go -destination=repo_mock_test.go -package book
type bookRepo interface {
	PaginateAllBooks(ctx context.Context, param book.PaginationParam) (book.PaginationResult, error)
	FindByIDs(ctx context.Context, id []string) ([]book.Book, error)
	PaginateBookSearch(ctx context.Context, searchQuery string, param book.PaginationParam) (book.PaginationResult, error)
}

// QueriesUseCase only act as a proxy.
type QueriesUseCase struct {
	repo bookRepo
}

func NewQueryUseCase(repo bookRepo) (*QueriesUseCase, error) {
	return &QueriesUseCase{repo: repo}, nil
}

func (q QueriesUseCase) GetAll(ctx context.Context, param book.PaginationParam) (book.PaginationResult, error) {
	return q.repo.PaginateAllBooks(ctx, param)
}
func (q QueriesUseCase) Search(ctx context.Context, searchQuery string, param book.PaginationParam) (book.PaginationResult, error) {
	return q.repo.PaginateBookSearch(ctx, searchQuery, param)
}

func (q QueriesUseCase) DetailByID(ctx context.Context, id string) (book.Book, error) {
	bookItem, err := q.repo.FindByIDs(ctx, []string{id})
	if err != nil {
		return book.Book{}, err
	}

	if len(bookItem) == 0 {
		return book.Book{}, bookrp.ErrNotFound
	}

	return bookItem[0], nil
}
