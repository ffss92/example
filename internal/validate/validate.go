package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// This is thread safe and should be a singleton for caching purposes.
// https://pkg.go.dev/github.com/go-playground/validator/v10@v10.16.0#New
var validate *validator.Validate

// Validates a struct using github.com/go-playground/validator tags.
func Struct(v any) error {
	return validate.Struct(v)
}

// Gets the instance of the validator.Validate.
func Validator() *validator.Validate {
	return validate
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
