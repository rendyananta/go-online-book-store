package book

import (
	"context"
	"net/http"

	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"github.com/rendyananta/example-online-book-store/internal/entity/pagination"
	apphttp "github.com/rendyananta/example-online-book-store/internal/presenter/http"
)

type queriesUseCase interface {
	GetAll(ctx context.Context, param book.PaginationParam) (book.PaginationResult, error)
	DetailByID(ctx context.Context, id string) (book.Book, error)
	Search(ctx context.Context, searchQuery string, param book.PaginationParam) (book.PaginationResult, error)
}

type Handler struct {
	Queries queriesUseCase
}

func (h Handler) Handle(server *http.ServeMux) {
	server.HandleFunc("GET /books", h.handleIndex)
	server.HandleFunc("GET /books/search", h.handleSearch)
	server.HandleFunc("GET /books/{id}", h.handleDetail)
}

func (h Handler) handleIndex(rw http.ResponseWriter, r *http.Request) {
	arw := apphttp.AppResponseWriter{}

	paginationResult, err := h.Queries.GetAll(r.Context(), book.PaginationParam{
		LastID: r.URL.Query().Get("last_id"),
	})

	if err != nil {
		arw.Write(rw, r, err)
		return
	}

	arw.Data = paginationResult.Data
	arw.Meta = pagination.PageInfo{
		PerPage: paginationResult.PerPage,
		LastID:  paginationResult.LastID,
	}
	arw.Write(rw, r, nil)
}

func (h Handler) handleDetail(rw http.ResponseWriter, r *http.Request) {
	var arw = &apphttp.AppResponseWriter{}

	item, err := h.Queries.DetailByID(r.Context(), r.PathValue("id"))

	if err != nil {
		arw.Write(rw, r, err)
		return
	}

	arw.Data = item
	arw.Write(rw, r, nil)
}

func (h Handler) handleSearch(rw http.ResponseWriter, r *http.Request) {
	var arw = &apphttp.AppResponseWriter{}

	query := r.URL.Query().Get("q")
	if query == "" {
		h.handleIndex(rw, r)
		return
	}

	paginationResult, err := h.Queries.Search(r.Context(), query, book.PaginationParam{
		LastID: r.URL.Query().Get("last_id"),
	})

	if err != nil {
		arw.Write(rw, r, err)
		return
	}

	arw.Data = paginationResult.Data
	arw.Meta = pagination.PageInfo{
		PerPage: paginationResult.PerPage,
		LastID:  paginationResult.LastID,
	}
	arw.Write(rw, r, nil)
}
