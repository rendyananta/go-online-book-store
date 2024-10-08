package order

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"github.com/rendyananta/example-online-book-store/internal/entity/order"
	"reflect"
	"testing"
)

func TestPlaceOrderUseCase_PlaceOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepoMock := NewMockorderRepo(ctrl)
	bookRepoMock := NewMockbookRepo(ctrl)

	type fields struct {
		orderRepo orderRepo
		bookRepo  bookRepo
	}
	type args struct {
		ctx   context.Context
		param order.Main
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       order.Main
		wantErr    bool
	}{
		{
			name: "can place order",
			fields: fields{
				orderRepo: orderRepoMock,
				bookRepo:  bookRepoMock,
			},
			args: args{
				ctx: context.Background(),
				param: order.Main{
					UserID: "1",
					Lines: []order.Line{
						{
							LineReferenceType: order.LineReferenceTypeBook,
							LineReferenceID:   "10",
							Quantity:          1,
						},
						{
							LineReferenceType: order.LineReferenceTypeBook,
							LineReferenceID:   "11",
							Quantity:          2,
						},
					}},
			},
			beforeTest: func() {
				bookRepoMock.EXPECT().FindByIDs(context.Background(), []string{"10", "11"}).
					Return([]book.Book{
						{
							ID:          "10",
							Title:       "Book 1",
							Description: "desc",
							Price:       1.2,
						},
						{
							ID:          "11",
							Title:       "Book 2",
							Description: "desc",
							Price:       1.2,
						},
					}, nil)

				orderInfo := order.Main{
					UserID:     "1",
					GrandTotal: 3.6,
					Status:     order.StatusDone,
					Lines: []order.Line{
						{
							LineReferenceType: order.LineReferenceTypeBook,
							LineReferenceID:   "10",
							LineItem: book.Book{
								ID:          "10",
								Title:       "Book 1",
								Description: "desc",
								Price:       1.2,
							},
							Quantity: 1,
							Amount:   1.2,
							Subtotal: 1.2,
						},
						{
							LineReferenceType: order.LineReferenceTypeBook,
							LineReferenceID:   "11",
							LineItem: book.Book{
								ID:          "11",
								Title:       "Book 2",
								Description: "desc",
								Price:       1.2,
							},
							Quantity: 2,
							Amount:   1.2,
							Subtotal: 2.4,
						},
					},
				}

				orderResult := order.Main{
					ID:         "1",
					UserID:     "1",
					GrandTotal: 3.6,
					Status:     order.StatusDone,
					Lines: []order.Line{
						{
							ID:                "1",
							OrderID:           "1",
							LineReferenceType: order.LineReferenceTypeBook,
							LineReferenceID:   "10",
							LineItem: book.Book{
								ID:          "10",
								Title:       "Book 1",
								Description: "desc",
								Price:       1.2,
							},
							Quantity: 1,
							Amount:   1.2,
							Subtotal: 1.2,
						},
						{
							ID:                "2",
							OrderID:           "1",
							LineReferenceType: order.LineReferenceTypeBook,
							LineReferenceID:   "11",
							LineItem: book.Book{
								ID:          "11",
								Title:       "Book 2",
								Description: "desc",
								Price:       1.2,
							},
							Quantity: 2,
							Amount:   1.2,
							Subtotal: 2.4,
						},
					},
				}

				orderRepoMock.EXPECT().Create(context.Background(), orderInfo).Return(orderResult, nil)
			},
			want: order.Main{
				ID:         "1",
				UserID:     "1",
				GrandTotal: 3.6,
				Status:     order.StatusDone,
				Lines: []order.Line{
					{
						ID:                "1",
						OrderID:           "1",
						LineReferenceType: order.LineReferenceTypeBook,
						LineReferenceID:   "10",
						LineItem: book.Book{
							ID:          "10",
							Title:       "Book 1",
							Description: "desc",
							Price:       1.2,
						},
						Quantity: 1,
						Amount:   1.2,
						Subtotal: 1.2,
					},
					{
						ID:                "2",
						OrderID:           "1",
						LineReferenceType: order.LineReferenceTypeBook,
						LineReferenceID:   "11",
						LineItem: book.Book{
							ID:          "11",
							Title:       "Book 2",
							Description: "desc",
							Price:       1.2,
						},
						Quantity: 2,
						Amount:   1.2,
						Subtotal: 2.4,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := PlaceOrderUseCase{
				orderRepo: tt.fields.orderRepo,
				bookRepo:  tt.fields.bookRepo,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := uc.PlaceOrder(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlaceOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
