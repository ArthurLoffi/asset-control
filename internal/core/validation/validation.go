package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"asset-control/internal/core/exception"

	"github.com/go-playground/validator/v10"
)

var engine = func() *validator.Validate {
	engine := validator.New()
	engine.SetTagName("binding")
	// Messages cite the public name (page_size), not the Go field (PageSize).
	engine.RegisterTagNameFunc(func(field reflect.StructField) string {
		for _, tag := range []string{"json", "form"} {
			if name := strings.Split(field.Tag.Get(tag), ",")[0]; name != "" && name != "-" {
				return name
			}
		}
		return field.Name
	})
	return engine
}()

// Struct validates `binding` tags and returns a *exception.ValidationError.
func Struct(target any) error {
	err := engine.Struct(target)
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return err
	}
	messages := make([]string, 0, len(fieldErrors))
	for _, fieldError := range fieldErrors {
		messages = append(messages, message(fieldError))
	}
	return exception.NewValidationError(messages...)
}

func message(fieldError validator.FieldError) string {
	field := fieldError.Field()
	switch fieldError.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "ip":
		return field + " must be a valid IP address"
	case "oneof":
		return fmt.Sprintf("%s must be one of the following values: %s", field, strings.ReplaceAll(fieldError.Param(), " ", ", "))
	default:
		return fmt.Sprintf("%s must satisfy %s=%s", field, fieldError.Tag(), fieldError.Param())
	}
}
