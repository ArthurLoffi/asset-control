package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"asset-control/internal/core/dto"
	"asset-control/internal/core/exception"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestActiveFilterIncludesLegacyRecordsAndPreservesCriteria(t *testing.T) {
	filter := bson.M{"employee_id": uint64(8)}
	active := activeFilter(filter)
	if active["employee_id"] != uint64(8) {
		t.Fatal("lost caller criteria")
	}
	if value, exists := active["deactivated_at"]; !exists || value != nil {
		t.Fatal("active reads must match missing or null timestamps")
	}
	if _, exists := filter["deactivated_at"]; exists {
		t.Fatal("modified caller filter")
	}
	if active := activeFilter(nil); len(active) != 1 {
		t.Fatalf("nil filter: %v", active)
	}
}

func TestUpdateFieldsPreserveExplicitNullAndFalseButNotLifecycle(t *testing.T) {
	type document struct {
		ID            uint64     `bson:"id"`
		HasMouse      *bool      `bson:"has_mouse"`
		EmployeeID    *uint64    `bson:"employee_id"`
		CreatedAt     time.Time  `bson:"created_at"`
		DeactivatedAt *time.Time `bson:"deactivated_at"`
	}
	isFalse := false
	now := time.Now()
	fields, err := documentFields(document{ID: 42, HasMouse: &isFalse, CreatedAt: now, DeactivatedAt: &now})
	if err != nil {
		t.Fatal(err)
	}
	for _, immutable := range []string{"_id", "id", "created_at", "deactivated_at"} {
		if _, exists := fields[immutable]; exists {
			t.Fatalf("immutable field %s included", immutable)
		}
	}
	if fields["has_mouse"] != false {
		t.Fatal("false value lost")
	}
	if value, exists := fields["employee_id"]; !exists || value != nil {
		t.Fatal("explicit null lost")
	}
}

func TestPersistenceErrors(t *testing.T) {
	duplicate := mongo.WriteException{WriteErrors: mongo.WriteErrors{{Code: 11000, Message: "duplicate"}}}
	for _, tc := range []struct{ input, expected error }{
		{nil, nil}, {mongo.ErrNoDocuments, exception.ErrNotFound}, {duplicate, exception.ErrConflict}, {context.Canceled, context.Canceled},
	} {
		if got := PersistenceError(tc.input); !errors.Is(got, tc.expected) {
			t.Fatalf("got %v, want %v", got, tc.expected)
		}
	}
	collection := Collection[struct{}]{notFound: func(id uint64) error { return exception.NotFoundByID("Machine", "MACHINE_001", id) }}
	if err := collection.persistenceError(mongo.ErrNoDocuments, 7); err == nil || err.Error() != "MACHINE_001: Machine not found by id" {
		t.Fatalf("got %v", err)
	}
}

func TestSortUsesAllowListAndIDTieBreaker(t *testing.T) {
	columns := map[string]string{"id": "id", "name": "name"}
	for _, tc := range []struct {
		request dto.PaginatedRequestDTO
		want    bson.D
	}{
		{dto.PaginatedRequestDTO{}, bson.D{{Key: "id", Value: 1}}},
		{dto.PaginatedRequestDTO{SortBy: "name", Direction: dto.DirectionDesc}, bson.D{{Key: "name", Value: -1}, {Key: "id", Value: -1}}},
		{dto.PaginatedRequestDTO{SortBy: "$where"}, bson.D{{Key: "id", Value: 1}}},
	} {
		got := Sort(tc.request, columns)
		if len(got) != len(tc.want) || got[0] != tc.want[0] {
			t.Errorf("%+v: got %v", tc.request, got)
		}
	}
	if Contains("a.b*")["$regex"] != `a\.b\*` {
		t.Fatal("search text must be literal")
	}
}
