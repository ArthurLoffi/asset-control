package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"asset-control/src/main/domain"
	"asset-control/src/main/infrastructure/outbound/data/entities"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestUpdateFieldsPreserveExplicitNullAndFalse(t *testing.T) {
	no := false
	fields, err := documentFields(entities.Machine{ID: 42, CreatedAt: time.Now(), HasMouse: &no})
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
		{nil, nil}, {mongo.ErrNoDocuments, domain.ErrNotFound}, {duplicate, domain.ErrConflict}, {context.Canceled, context.Canceled},
	} {
		if got := persistenceError(tt.input); !errors.Is(got, tt.expected) {
			t.Fatalf("got %v, want %v", got, tt.expected)
		}
	}
}
