package order

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendyananta/example-online-book-store/internal/entity/order"
	"github.com/rendyananta/example-online-book-store/pkg/db"
	"log/slog"
	"sync"
	"time"
)

//go:generate mockgen -source=order.go -destination=order_db_conn_mock_test.go -package order
type dbConnManager interface {
	Connection(name string) (*sqlx.DB, error)
}

type dbConnection interface {
	Preparex(query string) (*sqlx.Stmt, error)
	Rebind(query string) string
	Beginx() (*sqlx.Tx, error)
}

type preparedQueryGetter interface {
	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error
}

type Config struct {
	DBConn string
}

type preparedStmt struct {
	getUserOrdersPagination preparedQueryGetter
	getOrderDetail          preparedQueryGetter
	getOrderDetailByUser    preparedQueryGetter
	getOrderLines           preparedQueryGetter
}

type Repo struct {
	cfg          Config
	dbConn       dbConnection
	preparedStmt preparedStmt
}

func NewOrderRepo(cfg Config, dbConnManager dbConnManager) (*Repo, error) {
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

	r.preparedStmt.getUserOrdersPagination, err = r.dbConn.Preparex(r.dbConn.Rebind(queryGetUserOrdersPagination))
	if err != nil {
		return err
	}

	r.preparedStmt.getOrderDetail, err = r.dbConn.Preparex(r.dbConn.Rebind(queryGetOrderDetail))
	if err != nil {
		return err
	}

	r.preparedStmt.getOrderLines, err = r.dbConn.Preparex(r.dbConn.Rebind(queryGetOrderLines))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) PaginateOrdersByUserID(ctx context.Context, userID string, param order.PaginationParam) (order.PaginationResult, error) {
	if param.PerPage == 0 {
		param.PerPage = defaultPaginationLength
	}

	var result []tableOrder

	err := r.preparedStmt.getUserOrdersPagination.SelectContext(ctx, &result, userID, param.LastID, param.PerPage)
	if err != nil {
		return order.PaginationResult{PerPage: param.PerPage}, err
	}

	var orders = make([]order.Main, 0, len(result))

	for _, item := range result {
		var createdAt *time.Time
		var UpdatedAt *time.Time

		if item.CreatedAt.Valid {
			createdAt = &item.CreatedAt.Time
		}

		if item.UpdatedAt.Valid {
			UpdatedAt = &item.UpdatedAt.Time
		}

		orders = append(orders, order.Main{
			ID:         item.ID,
			UserID:     item.UserID,
			GrandTotal: item.GrandTotal,
			Status:     item.Status,
			CreatedAt:  createdAt,
			UpdatedAt:  UpdatedAt,
		})
	}

	if len(orders) == 0 {
		return order.PaginationResult{
			Data:    orders,
			PerPage: param.PerPage,
		}, nil
	}

	return order.PaginationResult{
		Data:    orders,
		PerPage: param.PerPage,
		LastID:  orders[len(orders)-1].ID,
	}, nil
}

func (r *Repo) GetDetailByID(ctx context.Context, orderID string) (order.Main, error) {
	var resMainOrder tableOrder
	var resOrderLines []tableOrderLine
	var errs []error

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := r.preparedStmt.getOrderDetail.GetContext(ctx, &resMainOrder, orderID)
		if err != nil {
			slog.Error("error get order detail query", slog.String("error", err.Error()), slog.String("order_id", orderID))
			errs = append(errs, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := r.preparedStmt.getOrderLines.SelectContext(ctx, &resOrderLines, orderID)
		if err != nil {
			slog.Error("error get order lines query", slog.String("error", err.Error()), slog.String("order_id", orderID))
			errs = append(errs, err)
		}
	}()

	wg.Wait()

	if len(errs) > 0 {
		return order.Main{}, errs[0]
	}

	orderLines := make([]order.Line, 0, len(resOrderLines))

	for _, line := range resOrderLines {
		orderLines = append(orderLines, order.Line{
			ID:                line.ID,
			OrderID:           line.OrderID,
			LineReferenceType: line.LineReferenceType,
			LineReferenceID:   line.LineReferenceID,
			Amount:            line.Amount,
			Quantity:          line.Quantity,
			Subtotal:          line.Subtotal,
		})
	}

	var createdAt *time.Time
	var UpdatedAt *time.Time

	if resMainOrder.CreatedAt.Valid {
		createdAt = &resMainOrder.CreatedAt.Time
	}

	if resMainOrder.UpdatedAt.Valid {
		UpdatedAt = &resMainOrder.UpdatedAt.Time
	}

	mainOrder := order.Main{
		ID:         resMainOrder.ID,
		UserID:     resMainOrder.UserID,
		GrandTotal: resMainOrder.GrandTotal,
		Status:     resMainOrder.Status,
		Lines:      orderLines,
		CreatedAt:  createdAt,
		UpdatedAt:  UpdatedAt,
	}

	return mainOrder, nil
}

func (r *Repo) Create(ctx context.Context, param order.Main) (order.Main, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return order.Main{}, err
	}

	var (
		createdAt time.Time
		updatedAt time.Time
	)

	if param.CreatedAt == nil {
		createdAt = time.Now()
	}

	if param.UpdatedAt == nil {
		updatedAt = time.Now()
	}

	tx, err := r.dbConn.Beginx()
	if err != nil {
		return order.Main{}, err
	}

	_, err = tx.ExecContext(ctx, queryInsertOrder, id.String(), param.UserID, param.GrandTotal, param.Status, createdAt, updatedAt)
	if err != nil {
		return order.Main{}, err
	}

	for _, line := range param.Lines {
		lineID, err := uuid.NewV7()
		if err != nil {
			continue
		}

		_, err = tx.ExecContext(ctx, queryInsertOrderLine, lineID.String(), id.String(), line.LineReferenceType, line.LineReferenceID, line.Amount, line.Quantity, line.Subtotal)
		if err != nil {
			slog.Error("error create order line", slog.String("error", err.Error()), slog.String("order_id", id.String()))
			continue
		}
	}

	err = tx.Commit()
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return order.Main{}, err
		}

		return order.Main{}, err
	}

	return order.Main{
		ID:         id.String(),
		UserID:     param.UserID,
		GrandTotal: param.GrandTotal,
		Status:     param.Status,
		Lines:      param.Lines,
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
	}, nil
}
