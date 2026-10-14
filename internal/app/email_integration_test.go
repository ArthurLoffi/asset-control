//go:build integration

package app

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"testing"
	"time"

	"asset-control/internal/core/database"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestEmailLookups(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	createdAt := time.Date(2025, 3, 4, 12, 0, 0, 0, time.UTC)
	collection := db.Collection(database.UsableEmailsCollection)
	router := testRouter(t, db)
	if listed, _ := page(t, callAPI(t, router, "GET", "/api/emails", nil, 200)); len(listed) != 0 {
		t.Fatalf("empty email list: %v", listed)
	}
	_, err := collection.InsertMany(ctx, []any{
		bson.M{"id": int32(43), "user": "Ana [Ops]", "email": "other@example.com", "password_hash": "other-password", "created_at": createdAt},
		bson.M{"id": int32(42), "user": "Ana [Ops]", "email": "ana@example.com", "password_hash": " plain<&>text password ", "created_at": createdAt},
		bson.M{"id": int32(1), "user": "Unrelated", "email": "unrelated@example.com", "password_hash": "unrelated", "created_at": createdAt},
	})
	if err != nil {
		t.Fatal(err)
	}
	listed, _ := page(t, callAPI(t, router, "GET", "/api/emails", nil, 200))
	if fmt.Sprint(ids(listed)) != "[1 42 43]" || listed[1]["password_hash"] != " plain<&>text password " {
		t.Fatal("email list did not preserve records and ID ordering")
	}
	data := object(t, callAPI(t, router, "GET", "/api/emails/42", nil, 200))
	if len(data) != 5 || data["id"] != float64(42) || data["user"] != "Ana [Ops]" || data["email"] != "ana@example.com" ||
		data["password_hash"] != " plain<&>text password " || data["created_at"] != createdAt.Format(time.RFC3339) {
		t.Fatal("lookup did not return the expected stored email")
	}
	if found, _ := page(t, callAPI(t, router, "GET", "/api/emails?user="+url.QueryEscape("[oPs]"), nil, 200)); fmt.Sprint(ids(found)) != "[42 43]" {
		t.Fatalf("literal case-insensitive user filter: %v", found)
	}
	for _, user := range []string{"missing", ".*", "ana@example.com"} {
		if found, _ := page(t, callAPI(t, router, "GET", "/api/emails?user="+url.QueryEscape(user), nil, 200)); len(found) != 0 {
			t.Fatalf("user %q matched %v", user, found)
		}
	}
	callAPI(t, router, "GET", "/api/emails/999", nil, 404)
	callAPI(t, router, "GET", "/api/emails/abc", nil, 400)
	var stored bson.M
	if err := collection.FindOne(ctx, bson.M{"id": int32(42)}).Decode(&stored); err != nil {
		t.Fatal(err)
	}
	if stored["password_hash"] != " plain<&>text password " || stored["id"] != int32(42) {
		t.Fatal("lookup changed the imported record")
	}
}

func TestEmailCreateAndUpdate(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	collection := db.Collection(database.UsableEmailsCollection)
	_, err := collection.InsertOne(ctx, bson.M{"id": int32(41), "user": "Imported", "email": "imported@example.com", "password_hash": "original", "created_at": time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	ensureSchema(t, db)
	router := testRouter(t, db)
	body := map[string]any{"user": "Ana", "email": "ana@example.com", "password_hash": " plain<&>text ", "id": 999, "created_at": "2000-01-01T00:00:00Z"}
	created := object(t, callAPI(t, router, "POST", "/api/emails", body, 201))
	if created["id"] != float64(42) || created["created_at"] == nil || created["created_at"] == body["created_at"] || created["password_hash"] != body["password_hash"] {
		t.Fatal("creation failed to generate fields or preserve password")
	}
	if _, err := collection.UpdateOne(ctx, bson.M{"id": int32(42)}, bson.M{"$set": bson.M{"migration_note": "keep"}}); err != nil {
		t.Fatal(err)
	}
	body["user"], body["email"], body["password_hash"] = "Updated", "updated@example.com", " changed password "
	updated := object(t, callAPI(t, router, "PUT", "/api/emails/42", body, 200))
	if updated["id"] != created["id"] || updated["created_at"] != created["created_at"] || updated["user"] != body["user"] ||
		updated["email"] != body["email"] || updated["password_hash"] != body["password_hash"] {
		t.Fatal("update failed to preserve immutable fields or replace mutable fields")
	}
	var stored bson.M
	if err := collection.FindOne(ctx, bson.M{"id": int32(42)}).Decode(&stored); err != nil {
		t.Fatal(err)
	}
	if stored["id"] != int32(42) || stored["migration_note"] != "keep" || stored["password_hash"] != body["password_hash"] {
		t.Fatal("update changed the ID type, unknown fields, or plaintext password")
	}
	callAPI(t, router, "PUT", "/api/emails/999", body, 404)
	callAPI(t, router, "POST", "/api/emails", map[string]any{}, 400)
	callAPI(t, router, "PUT", "/api/emails/42", map[string]any{}, 400)
	if count, err := collection.CountDocuments(ctx, bson.M{}); err != nil || count != 2 {
		t.Fatalf("unexpected document count: %d, %v", count, err)
	}
	delete(stored, "_id")
	if _, err := collection.InsertOne(ctx, stored); err == nil {
		t.Fatal("duplicate ID accepted")
	}
	// Exercise the final int32 ID and verify exhaustion does not advance the counter.
	_, err = db.Collection(database.CountersCollection).UpdateOne(ctx, bson.M{"_id": database.UsableEmailsCollection}, bson.M{"$set": bson.M{"seq": int64(math.MaxInt32 - 1)}})
	if err != nil {
		t.Fatal(err)
	}
	last := object(t, callAPI(t, router, "POST", "/api/emails", body, 201))
	if last["id"] != float64(math.MaxInt32) {
		t.Fatal("final int32 ID was not allocated")
	}
	callAPI(t, router, "POST", "/api/emails", body, 500)
	var counter struct {
		Seq int64 `bson:"seq"`
	}
	if err := db.Collection(database.CountersCollection).FindOne(ctx, bson.M{"_id": database.UsableEmailsCollection}).Decode(&counter); err != nil || counter.Seq != math.MaxInt32 {
		t.Fatalf("counter overflowed: %d, %v", counter.Seq, err)
	}
}
