//go:build integration

package app

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"asset-control/internal/modules/audit/domain/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestRecentChangesResolveEditorsAndDescribeUpdates(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	userID := bson.NewObjectID()
	if _, err := db.Collection("users").InsertOne(ctx, bson.M{
		"_id": userID, "name": "Ana", "password": "secret hash",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Collection("changes").InsertMany(ctx, []any{
		bson.M{
			"time":    time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC),
			"user_id": userID.Hex(), "resource": "machines", "resource_id": 42, "action": "update",
			"object": bson.M{"before": bson.M{"model": "Old"}, "after": bson.M{"model": "Middle"}},
		},
		bson.M{
			"time":    time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC),
			"user_id": userID.Hex(), "resource": "machines", "resource_id": 42, "action": "update",
			"object": bson.M{"before": bson.M{"model": "Middle", "password_hash": "[redacted]"},
				"after": bson.M{"model": "New", "password_hash": "[redacted]"}},
		},
	}); err != nil {
		t.Fatal(err)
	}
	var changes []dto.RecentChangeDto
	if err := json.Unmarshal(callAPI(t, testRouter(t, db), "GET", "/api/changes/recent", nil, 200), &changes); err != nil {
		t.Fatal(err)
	}
	if len(changes) != 2 || changes[0].UserName != "Ana" || changes[0].ChangedFields[0].After != "New" {
		t.Fatalf("unexpected recent changes: %#v", changes)
	}
	if len(changes[0].ChangedFields) != 1 || changes[0].ChangedFields[0].Field != "model" {
		t.Fatalf("secret or metadata in summary: %#v", changes[0].ChangedFields)
	}
}
