package order

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"github.com/rendyananta/example-online-book-store/internal/entity/order"
	"reflect"
	"testing"
)

func TestNewOrderQueriesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepoMock := NewMockorderRepo(ctrl)
	bookRepoMock := NewMockbookRepo(ctrl)

	type args struct {
		orderRepo orderRepo
		bookRepo  bookRepo
	}
	tests := []struct {
		name    string
		args    args
		want    *QueriesUseCase
		wantErr bool
	}{
		{
			name: "can init order queries",
			args: args{
				orderRepo: orderRepoMock,
				bookRepo:  bookRepoMock,
			},
			want: &QueriesUseCase{
				orderRepo: orderRepoMock,
				bookRepo:  bookRepoMock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrderQueriesUseCase(tt.args.orderRepo, tt.args.bookRepo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOrderQueriesUseCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderQueriesUseCase() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueriesUseCase_GetDetailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepoMock := NewMockorderRepo(ctrl)
	bookRepoMock := NewMockbookRepo(ctrl)

	type fields struct {
		orderRepo orderRepo
		bookRepo  bookRepo
	}
	type args struct {
		ctx     context.Context
		orderID string
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
			name: "can get detail by id",
			fields: fields{
				orderRepo: orderRepoMock,
				bookRepo:  bookRepoMock,
			},
			args: args{
				ctx:     context.Background(),
				orderID: "1",
			},
			beforeTest: func() {
				orderRepoMock.EXPECT().GetDetailByID(context.Background(), "1").
					Return(order.Main{
						ID:         "1",
						UserID:     "1",
						GrandTotal: 2.4,
						Status:     order.StatusDone,
						Lines: []order.Line{
							{
								ID:                "1",
								OrderID:           "1",
								LineReferenceType: order.LineReferenceTypeBook,
								LineReferenceID:   "10",
								Quantity:          1,
								Amount:            1.2,
								Subtotal:          1.2,
							},
							{
								ID:                "2",
								OrderID:           "1",
								LineReferenceType: order.LineReferenceTypeBook,
								LineReferenceID:   "11",
								Quantity:          1,
								Amount:            1.2,
								Subtotal:          1.2,
							},
						},
					}, nil)

				bookRepoMock.EXPECT().FindByIDs(context.Background(), []string{"10", "11"}).
					Return([]book.Book{
						{
							ID:          "10",
							Title:       "book 1",
							Description: "desc",
						},
						{
							ID:          "11",
							Title:       "book 2",
							Description: "desc",
						},
					}, nil)
			},
			want: order.Main{
				ID:         "1",
				UserID:     "1",
				GrandTotal: 2.4,
				Status:     order.StatusDone,
				Lines: []order.Line{
					{
						ID:                "1",
						OrderID:           "1",
						LineReferenceType: order.LineReferenceTypeBook,
						LineReferenceID:   "10",
						Quantity:          1,
						Amount:            1.2,
						Subtotal:          1.2,
						LineItem: book.Book{
							ID:          "10",
							Title:       "book 1",
							Description: "desc",
						},
					},
					{
						ID:                "2",
						OrderID:           "1",
						LineReferenceType: order.LineReferenceTypeBook,
						LineReferenceID:   "11",
						Quantity:          1,
						Amount:            1.2,
						Subtotal:          1.2,
						LineItem: book.Book{
							ID:          "11",
							Title:       "book 2",
							Description: "desc",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "can handle error get detail by id",
			fields: fields{
				orderRepo: orderRepoMock,
				bookRepo:  bookRepoMock,
			},
			args: args{
				ctx:     context.Background(),
				orderID: "1",
			},
			beforeTest: func() {
				orderRepoMock.EXPECT().GetDetailByID(context.Background(), "1").
					Return(order.Main{
						ID:         "1",
						UserID:     "1",
						GrandTotal: 2.4,
						Status:     order.StatusDone,
						Lines: []order.Line{
							{
								ID:                "1",
								OrderID:           "1",
								LineReferenceType: order.LineReferenceTypeBook,
								LineReferenceID:   "10",
								Quantity:          1,
								Amount:            1.2,
								Subtotal:          1.2,
							},
							{
								ID:                "2",
								OrderID:           "1",
								LineReferenceType: order.LineReferenceTypeBook,
								LineReferenceID:   "11",
								Quantity:          1,
								Amount:            1.2,
								Subtotal:          1.2,
							},
						},
					}, nil)

				bookRepoMock.EXPECT().FindByIDs(context.Background(), []string{"10", "11"}).
					Return(nil, sql.ErrConnDone)
			},
			want: order.Main{
				ID:         "1",
				UserID:     "1",
				GrandTotal: 2.4,
				Status:     order.StatusDone,
				Lines: []order.Line{
					{
						ID:                "1",
						OrderID:           "1",
						LineReferenceType: order.LineReferenceTypeBook,
						LineReferenceID:   "10",
						Quantity:          1,
						Amount:            1.2,
						Subtotal:          1.2,
					},
					{
						ID:                "2",
						OrderID:           "1",
						LineReferenceType: order.LineReferenceTypeBook,
						LineReferenceID:   "11",
						Quantity:          1,
						Amount:            1.2,
						Subtotal:          1.2,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := QueriesUseCase{
				orderRepo: tt.fields.orderRepo,
				bookRepo:  tt.fields.bookRepo,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := uc.GetDetailByID(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetailByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDetailByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueriesUseCase_PaginateOrdersByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepoMock := NewMockorderRepo(ctrl)
	bookRepoMock := NewMockbookRepo(ctrl)

	type fields struct {
		orderRepo orderRepo
		bookRepo  bookRepo
	}
	type args struct {
		ctx    context.Context
		userID string
		param  order.PaginationParam
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       order.PaginationResult
		wantErr    bool
	}{
		{
			name: "can paginate user orders",
			fields: fields{
				orderRepo: orderRepoMock,
				bookRepo:  bookRepoMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
				param: order.PaginationParam{
					PerPage: 2,
				},
			},
			beforeTest: func() {
				orderRepoMock.EXPECT().PaginateOrdersByUserID(context.Background(), "1", order.PaginationParam{
					PerPage: 2,
				}).Return(order.PaginationResult{
					Data: []order.Main{
						{
							ID:         "1",
							UserID:     "1",
							GrandTotal: 2.4,
							Status:     order.StatusDone,
						},
					},
					PerPage: 2,
					LastID:  "1",
				}, nil)
			},
			want: order.PaginationResult{
				Data: []order.Main{
					{
						ID:         "1",
						UserID:     "1",
						GrandTotal: 2.4,
						Status:     order.StatusDone,
					},
				},
				PerPage: 2,
				LastID:  "1",
			},
		},
		{
			name: "can handle error paginate user orders",
			fields: fields{
				orderRepo: orderRepoMock,
				bookRepo:  bookRepoMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
				param: order.PaginationParam{
					PerPage: 2,
				},
			},
			beforeTest: func() {
				orderRepoMock.EXPECT().PaginateOrdersByUserID(context.Background(), "1", order.PaginationParam{
					PerPage: 2,
				}).Return(order.PaginationResult{
					PerPage: 2,
					LastID:  "",
				}, sql.ErrConnDone)
			},
			want: order.PaginationResult{
				PerPage: 2,
				LastID:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := QueriesUseCase{
				orderRepo: tt.fields.orderRepo,
				bookRepo:  tt.fields.bookRepo,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := uc.PaginateOrdersByUserID(tt.args.ctx, tt.args.userID, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaginateOrdersByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginateOrdersByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
