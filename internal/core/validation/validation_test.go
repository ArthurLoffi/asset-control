package validation

import (
	"errors"
	"reflect"
	"testing"

	"asset-control/internal/core/exception"
)

func TestStructReportsPublicFieldNames(t *testing.T) {
	type input struct {
		Name      string   `json:"name" binding:"required"`
		Email     string   `json:"email" binding:"omitempty,email"`
		PageSize  int      `form:"page_size" binding:"lte=100"`
		Direction string   `form:"direction" binding:"oneof=ASC DESC"`
		Emails    []string `json:"emails" binding:"dive,max=3"`
	}
	err := Struct(input{Email: "invalid", PageSize: 101, Direction: "UP", Emails: []string{"long"}})
	var validationError *exception.ValidationError
	if !errors.As(err, &validationError) {
		t.Fatalf("got %T %v", err, err)
	}
	want := []string{
		"name is required",
		"email must be a valid email address",
		"page_size must satisfy lte=100",
		"direction must be one of the following values: ASC, DESC",
		"emails[0] must satisfy max=3",
	}
	if !reflect.DeepEqual(validationError.Messages, want) {
		t.Fatalf("got %q", validationError.Messages)
	}
	if err := Struct(input{Name: "Ana", Direction: "ASC"}); err != nil {
		t.Fatal(err)
	}
}
