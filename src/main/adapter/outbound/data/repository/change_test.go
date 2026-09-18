package repository

import (
	"testing"
	"time"

	"asset-control/src/main/domain/dto"
)

func TestAuditDocumentUsesJSONFieldNamesAndRedactsSecrets(t *testing.T) {
	input := dto.EmailDto{
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
