package order

import (
	"context"
	"math"

	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"github.com/rendyananta/example-online-book-store/internal/entity/order"
)

type PlaceOrderUseCase struct {
	orderRepo orderRepo
	bookRepo  bookRepo
}

func NewPlaceOrderUseCase(orderRepo orderRepo, bookRepo bookRepo) (*PlaceOrderUseCase, error) {
	return &PlaceOrderUseCase{
		orderRepo: orderRepo,
		bookRepo:  bookRepo,
	}, nil
}

func (uc PlaceOrderUseCase) PlaceOrder(ctx context.Context, param order.Main) (order.Main, error) {
	bookIDsToCheck := make([]string, 0)
	for _, line := range param.Lines {
		if line.LineReferenceType == order.LineReferenceTypeBook {
			bookIDsToCheck = append(bookIDsToCheck, line.LineReferenceID)
		}
	}

	books, err := uc.bookRepo.FindByIDs(ctx, bookIDsToCheck)
	if err != nil {
		return param, err
	}

	bookByID := make(map[string]book.Book)
	for _, bookItem := range books {
		bookByID[bookItem.ID] = bookItem
	}

	var grandTotal = 0.0
	for i, line := range param.Lines {
		// other than book type are ignored since its possibility only referencing to the book
		if line.LineReferenceType != order.LineReferenceTypeBook {
			continue
		}

		bookItem, bookExist := bookByID[line.LineReferenceID]
		if !bookExist {
			return param, ErrOrderLineInvalid
		}

		param.Lines[i].LineItem = bookItem
		param.Lines[i].Amount = bookItem.Price

		subtotalRounded := math.Ceil(bookItem.Price * float64(line.Quantity) * 100)

		// round to max two decimal
		param.Lines[i].Subtotal = subtotalRounded / 100
		grandTotal = math.Ceil(grandTotal*100+subtotalRounded) / 100
	}

	param.Status = order.StatusDone
	param.GrandTotal = grandTotal

	return uc.orderRepo.Create(ctx, param)
}
