//go:build integration

package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginAndInventoryIntegration(t *testing.T) {
	db := testDatabase(t)
	hash, err := bcrypt.GenerateFromPassword([]byte(" test password "), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	id := bson.NewObjectID()
	if _, err := db.Collection("users").InsertOne(context.Background(), bson.M{"_id": id, "name": "Auth Test", "password": string(hash)}); err != nil {
		t.Fatal(err)
	}
	router, authModule, err := newRouter(db, testAuthConfig())
	if err != nil {
		t.Fatal(err)
	}
	if err := authModule.EnsureSchema(context.Background()); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"user":"Auth Test","password":" test password "}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-CSRF-Protection", "1")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("%d %s", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly {
		t.Fatal("missing session cookie")
	}
	var audit bson.M
	if err := db.Collection("users").FindOne(context.Background(), bson.M{"_id": id}).Decode(&audit); err != nil {
		t.Fatal(err)
	}
	if audit["last_login_request_id"] != response.Header().Get("X-Request-ID") || audit["password"] != string(hash) {
		t.Fatal("invalid audit persistence")
	}
	for _, path := range []string{"/api/auth/session", "/api/emails"} {
		request = httptest.NewRequest("GET", path, nil)
		request.AddCookie(cookies[0])
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("%s: %d", path, response.Code)
		}
		var data map[string]any
		if err := json.Unmarshal(response.Body.Bytes(), &data); err != nil {
			t.Fatal(err)
		}
	}
	request = httptest.NewRequest("POST", "/api/auth/logout", nil)
	request.Header.Set("X-CSRF-Protection", "1")
	request.AddCookie(cookies[0])
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent || response.Result().Cookies()[0].MaxAge != -1 {
		t.Fatal("logout failed")
	}
	request = httptest.NewRequest("GET", "/api/emails", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != 401 {
		t.Fatal("inventory accessible without cookie")
	}
}
