package httpx

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/ArthurLoffi/asset-control-auth"
	"github.com/gin-gonic/gin"
)

func TestRecoveryDoesNotLogCredentials(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, nil)))
	defer slog.SetDefault(previous)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(auth.RequestID(), Recovery())
	r.GET("/panic", func(c *gin.Context) { panic("private password") })
	req := httptest.NewRequest("GET", "/panic", nil)
	req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: "private token"})
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	if response.Code != 500 || strings.Contains(output.String(), "private") || strings.Contains(response.Body.String(), "private") {
		t.Fatalf("unsafe recovery: %d %s", response.Code, output.String())
	}
	if response.Body.String() != `{"message":"Unexpected error","code":"CORE_003","logs":{}}` {
		t.Fatalf("unexpected body %s", response.Body.String())
	}
	if !strings.Contains(output.String(), response.Header().Get("X-Request-ID")) {
		t.Fatal("recovery log missing request ID")
	}
}
