package adapter

import (
	"reflect"
	"testing"
	"time"

	"asset-control/internal/modules/audit/domain/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestAuditDocumentUsesJSONFieldNamesAndRedactsSecrets(t *testing.T) {
	input := struct {
		ID           int32     `json:"id"`
		User         string    `json:"user"`
		Email        string    `json:"email"`
		PasswordHash string    `json:"password_hash"`
		CreatedAt    time.Time `json:"created_at"`
	}{
		ID:           42,
		User:         "arthur",
		Email:        "arthur@example.com",
		PasswordHash: "plain-secret",
		CreatedAt:    time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC),
	}

	document := auditDocument(input)
	if document["id"] != int32(42) {
		t.Fatalf("unexpected id: %#v", document["id"])
	}
	if document["password_hash"] != "[redacted]" {
		t.Fatalf("password was not redacted: %#v", document["password_hash"])
	}
	if _, exists := document["PasswordHash"]; exists {
		t.Fatal("Go field name leaked instead of JSON field name")
	}
	if _, ok := document["created_at"].(time.Time); !ok {
		t.Fatalf("created_at lost BSON date type: %T", document["created_at"])
	}
}

func TestChangedFieldsShowsValuesWithoutSecretsOrMetadata(t *testing.T) {
	before := bson.M{
		"id": int64(42), "created_at": time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		"model": "Old", "employee_id": nil, "password_hash": "[redacted]",
		"details": bson.M{"nested": "before"},
	}
	after := bson.M{
		"id": int64(42), "created_at": time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		"model": "New", "employee_id": int64(7), "password_hash": "different secret",
		"details": bson.M{"nested": "after"},
	}
	got := changedFields(before, after)
	want := []dto.ChangedFieldDto{
		{Field: "employee_id", Before: "", After: "7"},
		{Field: "model", Before: "Old", After: "New"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("changed fields: got %#v, want %#v", got, want)
	}
}
