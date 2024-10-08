package order

import (
	"context"
	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"github.com/rendyananta/example-online-book-store/internal/entity/order"
)

//go:generate mockgen -source=query.go -destination=repo_mock_test.go -package order
type orderRepo interface {
	PaginateOrdersByUserID(ctx context.Context, userID string, param order.PaginationParam) (order.PaginationResult, error)
	GetDetailByID(ctx context.Context, orderID string) (order.Main, error)
	Create(ctx context.Context, param order.Main) (order.Main, error)
}

type bookRepo interface {
	FindByIDs(ctx context.Context, id []string) ([]book.Book, error)
}

// QueriesUseCase only act as a proxy.
type QueriesUseCase struct {
	orderRepo orderRepo
	bookRepo  bookRepo
}

func NewOrderQueriesUseCase(orderRepo orderRepo, bookRepo bookRepo) (*QueriesUseCase, error) {
	return &QueriesUseCase{
		orderRepo: orderRepo,
		bookRepo:  bookRepo,
	}, nil
}

func (uc QueriesUseCase) PaginateOrdersByUserID(ctx context.Context, userID string, param order.PaginationParam) (order.PaginationResult, error) {
	return uc.orderRepo.PaginateOrdersByUserID(ctx, userID, param)
}

func (uc QueriesUseCase) GetDetailByID(ctx context.Context, orderID string) (order.Main, error) {
	orderDetail, err := uc.orderRepo.GetDetailByID(ctx, orderID)
	if err != nil {
		return orderDetail, err
	}

	lineIdxByBookID := make(map[string]int)
	bookIDs := make([]string, 0, len(orderDetail.Lines))
	for i, line := range orderDetail.Lines {
		if line.LineReferenceType == order.LineReferenceTypeBook {
			bookIDs = append(bookIDs, line.LineReferenceID)
			lineIdxByBookID[line.LineReferenceID] = i
		}
	}

	books, err := uc.bookRepo.FindByIDs(ctx, bookIDs)
	if err != nil {
		return orderDetail, err
	}

	for _, b := range books {
		lineIdx, ok := lineIdxByBookID[b.ID]
		if !ok {
			continue
		}

		orderDetail.Lines[lineIdx].LineItem = b
	}

	return orderDetail, nil
}
