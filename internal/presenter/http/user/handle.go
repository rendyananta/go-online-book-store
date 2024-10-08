package user

import (
	"context"
	"encoding/json"
	"errors"
	validatorpkg "github.com/go-playground/validator/v10"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	apphttp "github.com/rendyananta/example-online-book-store/internal/presenter/http"
	"github.com/rendyananta/example-online-book-store/pkg/validator"
	"net/http"
)

type registerUseCase interface {
	Register(ctx context.Context, param user.RegisterParam) (user.User, error)
	EmailRegistered(ctx context.Context, email string) bool
}

type authenticatorUseCase interface {
	Authenticate(ctx context.Context, param user.AuthenticateParam) (user.AuthenticateResult, error)
}

type Handler struct {
	Register      registerUseCase
	Authenticator authenticatorUseCase
}

func (h Handler) Handle(server *http.ServeMux) {
	server.HandleFunc("POST /auth/register", h.handleRegister)
	server.HandleFunc("POST /auth/token", h.handleToken)
}

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type AuthenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h Handler) handleRegister(rw http.ResponseWriter, r *http.Request) {
	var arw = &apphttp.AppResponseWriter{}
	var request RegisterRequest
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

	emailRegistered := h.Register.EmailRegistered(r.Context(), request.Email)
	if emailRegistered {
		arw.Write(rw, r, user.ErrEmailAlreadyRegistered)
		return
	}

	u, err := h.Register.Register(r.Context(), user.RegisterParam{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil || u.ID == "" {
		arw.Write(rw, r, err)
		return
	}

	arw.Data = u
	arw.Write(rw, r, nil)
	return
}

func (h Handler) handleToken(rw http.ResponseWriter, r *http.Request) {
	var arw = &apphttp.AppResponseWriter{}
	var request AuthenticateRequest
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

	ctx := r.Context()

	authenticateResult, err := h.Authenticator.Authenticate(ctx, user.AuthenticateParam{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		arw.Write(rw, r, err)
		return
	}

	arw.Data = authenticateResult
	arw.Write(rw, r, nil)
	return
}
