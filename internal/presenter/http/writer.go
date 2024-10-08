package http

import (
	"encoding/json"
	"errors"
	"github.com/rendyananta/example-online-book-store/pkg/validator"
	"log/slog"
	"net/http"

	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	"github.com/rendyananta/example-online-book-store/pkg/auth"
)

type Response struct {
	HTTPStatusCode int         `json:"-"`
	Message        string      `json:"message"`
	Data           interface{} `json:"data,omitempty"`
	Meta           interface{} `json:"meta,omitempty"`
	Errors         interface{} `json:"errors,omitempty"`
}

const defaultMessage = "OK"
const defaultStatusCode = http.StatusOK
const defaultInternalServerErrorResp = "{\"message\":\"Internal Server Error\"}"

var errWithResponse = map[error]Response{
	user.ErrEmailAlreadyRegistered: {
		Message:        "email already registered",
		HTTPStatusCode: http.StatusUnprocessableEntity,
	},
	user.ErrEmailIsNotRegistered: {
		Message:        "invalid credentials",
		HTTPStatusCode: http.StatusUnauthorized,
	},
	auth.ErrUnauthenticated: {
		Message:        "unauthenticated",
		HTTPStatusCode: http.StatusUnauthorized,
	},
}

type AppResponseWriter struct {
	StatusCode int
	Message    string
	Data       interface{}
	Meta       interface{}
}

func (d AppResponseWriter) Write(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		d.renderError(w, r, err)
		return
	}

	response := Response{
		Message: d.Message,
		Data:    d.Data,
		Meta:    d.Meta,
	}

	if response.Message == "" {
		response.Message = defaultMessage
	}

	bytesBuff, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(defaultInternalServerErrorResp))

		return
	}

	var statusCode = d.StatusCode
	if statusCode == 0 {
		statusCode = defaultStatusCode
	}

	w.WriteHeader(statusCode)
	w.Write(bytesBuff)
}

func (d AppResponseWriter) renderError(w http.ResponseWriter, r *http.Request, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errs = map[string][]string{}

		for _, validationError := range validationErrors {
			errs[validationError.Field()] = []string{validator.ErrMessage(validationError.Tag(), validationError.Field())}
		}

		resp := Response{
			HTTPStatusCode: http.StatusUnprocessableEntity,
			Message:        "unprocessable entity",
			Errors:         errs,
		}

		if bytesBuff, err := json.Marshal(resp); err == nil {
			w.Write(bytesBuff)
			return
		}
	}

	if resp, ok := errWithResponse[err]; ok {
		w.WriteHeader(resp.HTTPStatusCode)

		if bytesBuff, err := json.Marshal(resp); err == nil {
			w.Write(bytesBuff)
			return
		}
	}

	slog.Error("unable to complete the request", slog.String("error", err.Error()))

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(defaultInternalServerErrorResp))
}
