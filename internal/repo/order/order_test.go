package order

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/rendyananta/example-online-book-store/database/migrations"
	"github.com/rendyananta/example-online-book-store/internal/entity/order"
	"reflect"
	"testing"
)

func TestRepo_GetDetailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)
	preparedStmtMock := NewMockpreparedQueryGetter(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
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
			name: "can handle get orders",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					getUserOrdersPagination: preparedStmtMock,
					getOrderDetail:          preparedStmtMock,
					getOrderLines:           preparedStmtMock,
				},
			},
			args: args{
				ctx:     context.Background(),
				orderID: "1",
			},
			beforeTest: func() {
				var result tableOrder
				preparedStmtMock.EXPECT().GetContext(context.Background(), &result, "1").Return(nil).
					SetArg(1, tableOrder{
						ID:         "1",
						UserID:     "1",
						GrandTotal: 2.4,
						Status:     order.StatusDone,
						CreatedAt:  sql.NullTime{},
						UpdatedAt:  sql.NullTime{},
					})

				var orderLines []tableOrderLine
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &orderLines, "1").Return(nil).
					SetArg(1, []tableOrderLine{
						{
							ID:                "1",
							OrderID:           "1",
							LineReferenceType: "book",
							LineReferenceID:   "1",
							Amount:            4.5,
							Quantity:          2,
							Subtotal:          9,
						},
					})
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
						LineReferenceType: "book",
						LineReferenceID:   "1",
						Amount:            4.5,
						Quantity:          2,
						Subtotal:          9,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "can handle get orders",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					getUserOrdersPagination: preparedStmtMock,
					getOrderDetail:          preparedStmtMock,
					getOrderLines:           preparedStmtMock,
				},
			},
			args: args{
				ctx:     context.Background(),
				orderID: "1",
			},
			beforeTest: func() {
				var result tableOrder
				preparedStmtMock.EXPECT().GetContext(context.Background(), &result, "1").Return(nil).
					SetArg(1, tableOrder{
						ID:         "1",
						UserID:     "1",
						GrandTotal: 2.4,
						Status:     order.StatusDone,
						CreatedAt:  sql.NullTime{},
						UpdatedAt:  sql.NullTime{},
					})

				var orderLines []tableOrderLine
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &orderLines, "1").Return(sql.ErrConnDone)
			},
			want:    order.Main{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := r.GetDetailByID(tt.args.ctx, tt.args.orderID)
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

func TestRepo_PaginateOrdersByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)
	preparedStmtMock := NewMockpreparedQueryGetter(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
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
			name: "can handle get orders",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					getUserOrdersPagination: preparedStmtMock,
				},
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
				param: order.PaginationParam{
					PerPage: 2,
				},
			},
			beforeTest: func() {
				var result []tableOrder
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &result, "1", "", 2).Return(nil).
					SetArg(1, []tableOrder{
						{
							ID:         "1",
							UserID:     "1",
							GrandTotal: 2.4,
							Status:     order.StatusDone,
							CreatedAt:  sql.NullTime{},
							UpdatedAt:  sql.NullTime{},
						},
						{
							ID:         "2",
							UserID:     "1",
							GrandTotal: 3.4,
							Status:     order.StatusDone,
							CreatedAt:  sql.NullTime{},
							UpdatedAt:  sql.NullTime{},
						},
					})
			},
			want: order.PaginationResult{
				Data: []order.Main{
					{
						ID:         "1",
						UserID:     "1",
						GrandTotal: 2.4,
						Status:     order.StatusDone,
					},
					{
						ID:         "2",
						UserID:     "1",
						GrandTotal: 3.4,
						Status:     order.StatusDone,
					},
				},
				PerPage: 2,
				LastID:  "2",
			},
			wantErr: false,
		},
		{
			name: "can handle get orders",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					getUserOrdersPagination: preparedStmtMock,
				},
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
				param: order.PaginationParam{
					PerPage: 2,
				},
			},
			beforeTest: func() {
				var result []tableOrder
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &result, "1", "", 2).Return(sql.ErrConnDone)
			},
			want: order.PaginationResult{
				Data:    nil,
				PerPage: 2,
				LastID:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := r.PaginateOrdersByUserID(tt.args.ctx, tt.args.userID, tt.args.param)
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

func TestRepo_boot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	tests := []struct {
		name       string
		fields     fields
		beforeTest func()
		wantErr    bool
	}{
		{
			name: "success boot repo",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			beforeTest: func() {
				dbConnMock.EXPECT().Rebind(queryGetUserOrdersPagination).Return(queryGetUserOrdersPagination)
				dbConnMock.EXPECT().Preparex(queryGetUserOrdersPagination).Return(&sqlx.Stmt{}, nil)

				dbConnMock.EXPECT().Rebind(queryGetOrderDetail).Return(queryGetOrderDetail)
				dbConnMock.EXPECT().Preparex(queryGetOrderDetail).Return(&sqlx.Stmt{}, nil)

				dbConnMock.EXPECT().Rebind(queryGetOrderLines).Return(queryGetOrderLines)
				dbConnMock.EXPECT().Preparex(queryGetOrderLines).Return(&sqlx.Stmt{}, nil)
			},
			wantErr: false,
		},
		{
			name: "can handle error boot repo",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			beforeTest: func() {
				dbConnMock.EXPECT().Rebind(queryGetUserOrdersPagination).Return(queryGetUserOrdersPagination)
				dbConnMock.EXPECT().Preparex(queryGetUserOrdersPagination).Return(&sqlx.Stmt{}, nil)

				dbConnMock.EXPECT().Rebind(queryGetOrderDetail).Return(queryGetOrderDetail)
				dbConnMock.EXPECT().Preparex(queryGetOrderDetail).Return(&sqlx.Stmt{}, nil)

				dbConnMock.EXPECT().Rebind(queryGetOrderLines).Return(queryGetOrderLines)
				dbConnMock.EXPECT().Preparex(queryGetOrderLines).Return(&sqlx.Stmt{}, sql.ErrConnDone)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}

			if tt.beforeTest != nil {
				tt.beforeTest()
			}

			if err := r.boot(); (err != nil) != tt.wantErr {
				t.Errorf("boot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock, _ := sqlx.Open("sqlite3", ":memory:")
	_ = migrations.CreateOrdersTable{Conn: dbConnMock}.Up()
	_ = migrations.CreateOrderLinesTable{Conn: dbConnMock}.Up()

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx   context.Context
		param order.Main
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(repo *Repo)
		want       order.Main
		wantErr    bool
	}{
		{
			name: "can create order",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			args: args{
				ctx: context.Background(),
				param: order.Main{
					UserID:     "1",
					GrandTotal: 2.4,
					Status:     order.StatusDone,
					Lines: []order.Line{
						{
							LineReferenceType: order.LineReferenceTypeBook,
							LineReferenceID:   "2",
							Amount:            1.2,
							Quantity:          2,
							Subtotal:          2.4,
						},
					},
				},
			},
			beforeTest: func(repo *Repo) {
			},
			want: order.Main{
				ID:         "",
				UserID:     "1",
				GrandTotal: 2.4,
				Status:     order.StatusDone,
				Lines: []order.Line{
					{
						LineReferenceType: order.LineReferenceTypeBook,
						LineReferenceID:   "2",
						Amount:            1.2,
						Quantity:          2,
						Subtotal:          2.4,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}
			if tt.beforeTest != nil {
				tt.beforeTest(r)
			}
			got, err := r.Create(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// unset the id
			got.ID = ""
			got.CreatedAt = nil
			got.UpdatedAt = nil

			for i := range got.Lines {
				got.Lines[i].ID = ""
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
