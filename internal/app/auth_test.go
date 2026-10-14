package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	auth "github.com/Brasil-Novo/asset-control-auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func testAuthConfig() auth.Config {
	return auth.Config{Secret: strings.Repeat("test-only-key-", 3), TTL: time.Hour, IsSecure: false, Origins: []string{"http://localhost:5173"}}
}

func newRouter(db *mongo.Database, cfg auth.Config) (*gin.Engine, *auth.Module, error) {
	authModule, err := auth.New(db, cfg)
	if err != nil {
		return nil, nil, err
	}
	router, err := New(db, authModule, nil)
	return router, authModule, err
}

func TestTunnelClientIPOnlyFromLocalProxy(t *testing.T) {
	r := gin.New()
	if err := configureClientIP(r); err != nil {
		t.Fatal(err)
	}
	r.GET("/ip", func(c *gin.Context) { c.String(http.StatusOK, c.ClientIP()) })
	for _, tc := range []struct {
		name, remote, cloudflareIP, forwardedFor, want string
	}{
		{"cloudflared IPv6", "[::1]:1234", "203.0.113.10", "198.51.100.1", "203.0.113.10"},
		{"cloudflared IPv4", "127.0.0.1:1234", "2001:db8::10", "198.51.100.1", "2001:db8::10"},
		{"direct client", "192.0.2.10:1234", "203.0.113.10", "198.51.100.1", "192.0.2.10"},
		{"missing Cloudflare IP", "[::1]:1234", "", "198.51.100.1", "::1"},
		{"invalid Cloudflare IP", "[::1]:1234", "not-an-ip", "198.51.100.1", "::1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ip", nil)
			req.RemoteAddr = tc.remote
			req.Header.Set(gin.PlatformCloudflare, tc.cloudflareIP)
			req.Header.Set("X-Forwarded-For", tc.forwardedFor)
			response := httptest.NewRecorder()
			r.ServeHTTP(response, req)
			if response.Code != http.StatusOK || response.Body.String() != tc.want {
				t.Fatalf("status=%d client_ip=%q, want %q", response.Code, response.Body.String(), tc.want)
			}
		})
	}
}

func TestSessionBehindTrustedTunnel(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	router, authModule, err := newRouter(client.Database("unused"), testAuthConfig())
	if err != nil {
		t.Fatal(err)
	}
	cookie, err := authModule.SessionCookie("test-user", "Test User")
	if err != nil {
		t.Fatal(err)
	}
	raw := cookie.Value
	for _, tc := range []struct {
		name, cookie string
		want         int
	}{
		{"valid cookie", raw, http.StatusOK},
		{"missing cookie", "", http.StatusUnauthorized},
		{"invalid cookie", "forged", http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/auth/session", nil)
			req.RemoteAddr = "[::1]:1234"
			req.Header.Set(gin.PlatformCloudflare, "203.0.113.10")
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: tc.cookie})
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tc.want {
				t.Fatalf("status=%d, want %d: %s", response.Code, tc.want, response.Body.String())
			}
		})
	}
}

func testRouter(t *testing.T, db *mongo.Database) http.Handler {
	t.Helper()
	cfg := testAuthConfig()
	router, authModule, err := newRouter(db, cfg)
	if err != nil {
		t.Fatal(err)
	}
	cookie, err := authModule.SessionCookie("test-user", "Test User")
	if err != nil {
		t.Fatal(err)
	}
	var requestCount atomic.Int64
	// Existing inventory tests still exercise real token validation. Each request
	// comes from a distinct client so suites never exhaust the per-IP rate limit.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		r.RemoteAddr = fmt.Sprintf("198.18.%d.%d:1234", count/250%250, count%250+1)
		r.AddCookie(cookie)
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
	router, _, err := newRouter(client.Database("unused"), testAuthConfig())
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
				req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: cookie})
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
	if response.Code != http.StatusNoContent {
		t.Fatal(response.Code)
	}
}

func TestUnknownRoutesUseErrorEnvelope(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	response := httptest.NewRecorder()
	testRouter(t, client.Database("unused")).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/unknown", nil))
	if response.Code != http.StatusNotFound || response.Body.String() != `{"message":"Cannot GET /api/unknown","code":"CORE_002","logs":{}}` {
		t.Fatalf("%d %s", response.Code, response.Body.String())
	}
}
