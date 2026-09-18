package main

import (
	"context"
	"fmt"
	"github.com/Brasil-Novo/asset-control-auth/src/main/adapter/security"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testAuthConfig() security.Config {
	return security.Config{Secret: strings.Repeat("test-only-key-", 3), TTL: time.Hour, Secure: false, Origins: []string{"http://localhost:5173"}}
}
func testRouter(t *testing.T, db *mongo.Database) http.Handler {
	t.Helper()
	cfg := testAuthConfig()
	router, err := newRouter(db, cfg)
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := security.NewTokens(cfg.Secret, cfg.TTL)
	if err != nil {
		t.Fatal(err)
	}
	raw, _, err := tokens.Generate("test-user", "Test User")
	if err != nil {
		t.Fatal(err)
	}
	// Existing inventory tests still exercise real token validation.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.AddCookie(&http.Cookie{Name: security.CookieName, Value: raw})
		r.Header.Set("X-CSRF-Protection", "1")
		router.ServeHTTP(w, r)
	})
}
func TestAllInventoryRoutesRequireAuthentication(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	router, err := newRouter(client.Database("unused"), testAuthConfig())
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, route := range router.Routes() {
		if strings.HasPrefix(route.Path, "/api/auth/") {
			continue
		}
		path := strings.ReplaceAll(route.Path, "*any", "index.html")
		for _, cookie := range []string{"", "forged"} {
			req := httptest.NewRequest(route.Method, path, nil)
			req.RemoteAddr = fmt.Sprintf("198.51.100.%d:12345", count+1)
			req.Header.Set("X-CSRF-Protection", "1")
			if cookie != "" {
				req.AddCookie(&http.Cookie{Name: security.CookieName, Value: cookie})
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != 401 {
				t.Fatalf("%s %s: %d", route.Method, path, response.Code)
			}
		}
		count++
	}
	if count < 30 {
		t.Fatalf("unexpected route coverage: %d", count)
	}
	// Logout must remain reachable with no session.
	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	req.RemoteAddr = "198.51.100.254:12345"
	req.Header.Set("X-CSRF-Protection", "1")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if response.Code != 200 {
		t.Fatal(response.Code)
	}
}
func TestInventoryStartupRequiresAuthConfiguration(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("AUTH_JWT_SECRET", "")
	if err := run(); err == nil {
		t.Fatal("started without signing key")
	}
}
