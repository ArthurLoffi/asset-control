package adapter

import (
	"asset-control/internal/core/database"
	"context"
	"errors"
	"testing"
	"time"

	exception "asset-control/internal/core/exception"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestUpdateFieldsPreserveExplicitNullAndFalse(t *testing.T) {
	no := false
	fields, err := documentFields(schema.Machine{ID: 42, CreatedAt: time.Now(), HasMouse: &no})
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"_id", "id", "created_at"} {
		if _, exists := fields[key]; exists {
			t.Fatalf("immutable field %s included", key)
		}
	}
	if fields["has_mouse"] != false {
		t.Fatal("false value lost")
	}
	if value, exists := fields["employee_id"]; !exists || value != nil {
		t.Fatal("explicit null lost")
	}
	if _, exists := fields["serial_number"]; !exists {
		t.Fatal("snake_case mapping lost")
	}
}

func TestPersistenceErrors(t *testing.T) {
	duplicate := mongo.WriteException{WriteErrors: mongo.WriteErrors{{Code: 11000, Message: "duplicate"}}}
	for _, tt := range []struct{ input, expected error }{
		{nil, nil}, {mongo.ErrNoDocuments, exception.ErrNotFound}, {duplicate, exception.ErrConflict}, {context.Canceled, context.Canceled},
	} {
		if got := database.PersistenceError(tt.input); !errors.Is(got, tt.expected) {
			t.Fatalf("got %v, want %v", got, tt.expected)
		}
	}
}
