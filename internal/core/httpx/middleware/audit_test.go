package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"asset-control/internal/core/audit"

	auth "github.com/ArthurLoffi/asset-control-auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestAuditIdentityCopiesAuthenticatedUserToRequestContext(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	authModule, err := auth.New(client.Database("unused"), auth.Config{
		Secret: strings.Repeat("k", 32), TTL: time.Hour, Origins: []string{"http://localhost:5173"},
	})
	if err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(authModule.RequireSession(), AuditIdentity())
	router.GET("/", func(c *gin.Context) {
		userID, ok := audit.UserID(c.Request.Context())
		if !ok || userID != "user-id" {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})
	cookie, err := authModule.SessionCookie("user-id", "Arthur")
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("got status %d", recorder.Code)
	}
}
