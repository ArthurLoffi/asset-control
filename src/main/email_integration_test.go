//go:build integration

package main

import (
	"context"
	"encoding/json"
	"math"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"asset-control/src/main/infrastructure/outbound/data/database"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestEmailLookups(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	createdAt := time.Date(2025, 3, 4, 12, 0, 0, 0, time.UTC)
	collection := db.Collection(database.UsableEmailsCollection)
	router := testRouter(t, db)
	if data := callAPI(t, router, "GET", "/api/emails", nil, 200); string(data) != "[]" {
		t.Fatalf("empty email list: %s", data)
	}
	_, err := collection.InsertMany(ctx, []any{
		bson.M{"id": int32(43), "user": "Ana [Ops]", "email": "other@example.com", "password_hash": "other-password", "created_at": createdAt},
		bson.M{"id": int32(42), "user": "Ana [Ops]", "email": "ana@example.com", "password_hash": " plain<&>text password ", "created_at": createdAt},
		bson.M{"id": int32(1), "user": "Unrelated", "email": "unrelated@example.com", "password_hash": "unrelated", "created_at": createdAt},
	})
	if err != nil {
		t.Fatal(err)
	}
	var listed []map[string]any
	if err := json.Unmarshal(callAPI(t, router, "GET", "/api/emails", nil, 200), &listed); err != nil {
		t.Fatal(err)
	}
	if len(listed) != 3 || listed[0]["id"] != float64(1) || listed[1]["id"] != float64(42) || listed[2]["id"] != float64(43) || listed[1]["password_hash"] != " plain<&>text password " {
		t.Fatal("email list did not preserve records and ID ordering")
	}
	for _, target := range []string{"/api/emails/id?id=42", "/api/emails/user?user=" + url.QueryEscape(" [oPs] ")} {
		data := object(t, callAPI(t, router, "GET", target, nil, 200))
		if len(data) != 5 || data["id"] != float64(42) || data["user"] != "Ana [Ops]" || data["email"] != "ana@example.com" || data["password_hash"] != " plain<&>text password " || data["created_at"] != createdAt.Format(time.RFC3339) {
			t.Fatal("lookup did not return the expected stored email")
		}
	}
	for _, target := range []string{"/api/emails/id?id=999", "/api/emails/user?user=missing", "/api/emails/user?user=.*", "/api/emails/user?user=ana@example.com"} {
		callAPI(t, router, "GET", target, nil, 404)
	}
	for _, target := range []string{"/api/emails/id", "/api/emails/id?id=2147483648", "/api/emails/user?user=%20"} {
		callAPI(t, router, "GET", target, nil, 400)
	}
	for _, method := range []string{"DELETE"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(method, "/api/emails", nil))
		if recorder.Code != 404 {
			t.Fatalf("unexpected email route: %s /api/emails returned %d", method, recorder.Code)
		}
	}
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
	if err := database.EnsureSchema(ctx, db); err != nil {
		t.Fatal(err)
	}
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
	updated := object(t, callAPI(t, router, "PUT", "/api/emails?id=42", body, 200))
	if updated["id"] != created["id"] || updated["created_at"] != created["created_at"] || updated["user"] != body["user"] || updated["email"] != body["email"] || updated["password_hash"] != body["password_hash"] {
		t.Fatal("update failed to preserve immutable fields or replace mutable fields")
	}
	var stored bson.M
	if err := collection.FindOne(ctx, bson.M{"id": int32(42)}).Decode(&stored); err != nil {
		t.Fatal(err)
	}
	if stored["id"] != int32(42) || stored["migration_note"] != "keep" || stored["password_hash"] != body["password_hash"] {
		t.Fatal("update changed the ID type, unknown fields, or plaintext password")
	}
	callAPI(t, router, "PUT", "/api/emails?id=999", body, 404)
	callAPI(t, router, "POST", "/api/emails", map[string]any{}, 400)
	callAPI(t, router, "PUT", "/api/emails?id=42", map[string]any{}, 400)
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
