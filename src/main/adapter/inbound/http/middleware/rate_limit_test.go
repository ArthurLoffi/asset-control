package middleware

import (
	"net/http"
	"net/http/httptest"
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

	// O burst permite 5 requisições imediatamente.
	for i := 1; i <= 5; i++ {
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

	// A 6ª requisição deve ser bloqueada.
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
