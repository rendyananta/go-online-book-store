package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type ValidationErrors = validator.ValidationErrors

var v = validator.New(validator.WithRequiredStructEnabled())

func SetUp() {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})
}

func Struct(any any) error {
	return v.Struct(any)
}
