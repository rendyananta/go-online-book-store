package order

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	validatorpkg "github.com/go-playground/validator/v10"
	httpen "github.com/rendyananta/example-online-book-store/internal/entity/http"
	"github.com/rendyananta/example-online-book-store/internal/entity/order"
	apphttp "github.com/rendyananta/example-online-book-store/internal/presenter/http"
	"github.com/rendyananta/example-online-book-store/pkg/auth"
	"github.com/rendyananta/example-online-book-store/pkg/validator"
)

type authMiddleware interface {
	Handle(next http.Handler) http.Handler
}

type placeOrderUseCase interface {
	PlaceOrder(ctx context.Context, param order.Main) (order.Main, error)
}

type queriesUseCase interface {
	PaginateOrdersByUserID(ctx context.Context, userID string, param order.PaginationParam) (order.PaginationResult, error)
	GetDetailByID(ctx context.Context, orderID string) (order.Main, error)
}

type Handler struct {
	AuthMiddleware    authMiddleware
	PlaceOrderUseCase placeOrderUseCase
	Queries           queriesUseCase
}

func (h Handler) Handle(server *http.ServeMux) {
	server.Handle("GET /orders", h.AuthMiddleware.Handle(http.HandlerFunc(h.handleIndex)))
	server.Handle("POST /orders/place", h.AuthMiddleware.Handle(http.HandlerFunc(h.handlePlaceOrder)))
	server.Handle("GET /orders/{id}", h.AuthMiddleware.Handle(http.HandlerFunc(h.handleDetail)))
}

func (h Handler) handleIndex(rw http.ResponseWriter, r *http.Request) {
	var arw = &apphttp.AppResponseWriter{}
	ctx := r.Context()
	userSession := ctx.Value(auth.CtxKeyUserSession).(*auth.UserSession)
	if userSession == nil {
		arw.Write(rw, r, auth.ErrUnauthenticated)
		return
	}

	result, err := h.Queries.PaginateOrdersByUserID(ctx, userSession.ID, order.PaginationParam{
		LastID: r.URL.Query().Get("last_id"),
	})

	if err != nil {
		arw.Write(rw, r, auth.ErrUnauthenticated)
		return
	}

	arw.Data = result

	arw.Write(rw, r, nil)
}

func (h Handler) handleDetail(rw http.ResponseWriter, r *http.Request) {
	var arw = &apphttp.AppResponseWriter{}
	ctx := r.Context()
	userSession := ctx.Value(auth.CtxKeyUserSession).(*auth.UserSession)

	item, err := h.Queries.GetDetailByID(r.Context(), r.PathValue("id"))
	if err != nil {
		arw.Write(rw, r, err)
		return
	}

	if item.UserID != userSession.ID {
		arw.Write(rw, r, httpen.ErrUnauthorized)
	}

	arw.Data = item
	arw.Write(rw, r, nil)
}

type LineItem struct {
	LineReferenceType string `json:"line_reference_type" validate:"required,oneof=book"`
	LineReferenceID   string `json:"line_reference_id" validate:"required"`
	Quantity          int    `json:"quantity" validate:"required,gt=0"`
}

type PlaceOrderRequest struct {
	Lines []LineItem `json:"lines" validate:"gt=0,dive"`
}

func (h Handler) handlePlaceOrder(rw http.ResponseWriter, r *http.Request) {
	var arw = &apphttp.AppResponseWriter{}
	ctx := r.Context()
	userSession := ctx.Value(auth.CtxKeyUserSession).(*auth.UserSession)
	if userSession == nil {
		arw.Write(rw, r, auth.ErrUnauthenticated)
		return
	}

	var request PlaceOrderRequest
	var err error

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			arw.Write(rw, r, err)
			return
		}
	}

	err = validator.Struct(request)
	var validationErrors validatorpkg.ValidationErrors
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		arw.Write(rw, r, err)
		return
	}

	orderLines := make([]order.Line, 0, len(request.Lines))
	for _, line := range request.Lines {
		orderLines = append(orderLines, order.Line{
			LineReferenceType: line.LineReferenceType,
			LineReferenceID:   line.LineReferenceID,
			Quantity:          line.Quantity,
		})
	}

	orderParam := order.Main{
		UserID: userSession.ID,
		Lines:  orderLines,
	}

	orderDetail, err := h.PlaceOrderUseCase.PlaceOrder(ctx, orderParam)
	if err != nil {
		arw.Write(rw, r, err)
		return
	}

	arw.Data = orderDetail
	arw.Write(rw, r, nil)
}
