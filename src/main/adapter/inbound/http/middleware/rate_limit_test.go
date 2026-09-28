package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Garante que o teste começa sem IPs já registrados.
	mu.Lock()
	clients = make(map[string]*client)
	mu.Unlock()

	router := gin.New()
	_ = router.SetTrustedProxies(nil)

	router.Use(RateLimit())

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	const ip = "192.168.1.10:12345"

	// O burst permite 15 requisições imediatamente.
	for i := 1; i <= 15; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = ip

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"requisição %d: esperado status %d, recebido %d",
				i,
				http.StatusOK,
				rec.Code,
			)
		}
	}

	// A 16ª requisição deve ser bloqueada.
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"esperado status %d, recebido %d",
			http.StatusTooManyRequests,
			rec.Code,
		)
	}

	expectedBody := `{"error":"Too many requests. Try again later."}`

	if rec.Body.String() != expectedBody {
		t.Errorf(
			"resposta esperada %q, recebida %q",
			expectedBody,
			rec.Body.String(),
		)
	}
}

func TestRateLimitUsesTrustedClientIP(t *testing.T) {
	router := gin.New()
	router.RemoteIPHeaders = []string{gin.PlatformCloudflare}
	if err := router.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		t.Fatal(err)
	}
	router.Use(RateLimit())
	router.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	request := func(remote, forwarded string) int {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = remote
		req.Header.Set(gin.PlatformCloudflare, forwarded)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		return response.Code
	}
	for i := 0; i < 15; i++ {
		if status := request("[::1]:1234", "203.0.113.21"); status != http.StatusOK {
			t.Fatalf("first visitor request %d: %d", i, status)
		}
	}
	if status := request("[::1]:1234", "203.0.113.21"); status != http.StatusTooManyRequests {
		t.Fatalf("first visitor exceeded burst: %d", status)
	}
	if status := request("[::1]:1234", "203.0.113.22"); status != http.StatusOK {
		t.Fatalf("second visitor shared the limiter: %d", status)
	}
	for i := 0; i < 15; i++ {
		forwarded := "203.0.113." + strconv.Itoa(i+30)
		if status := request("192.0.2.21:1234", forwarded); status != http.StatusOK {
			t.Fatalf("direct visitor request %d: %d", i, status)
		}
	}
	if status := request("192.0.2.21:1234", "203.0.113.45"); status != http.StatusTooManyRequests {
		t.Fatalf("direct visitor spoofed forwarded IP: %d", status)
	}
}
