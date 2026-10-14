package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"asset-control/src/main/domain"

	authMiddleware "github.com/Brasil-Novo/asset-control-auth/src/main/adapter/inbound/http/middleware"
	authDto "github.com/Brasil-Novo/asset-control-auth/src/main/domain/dto"
	"github.com/gin-gonic/gin"
)

func TestAuditIdentityCopiesAuthenticatedUserToRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(authMiddleware.SessionKey, authDto.SessionDto{ID: "user-id", Name: "Arthur"})
		c.Next()
	})
	router.Use(AuditIdentity())
	router.GET("/", func(c *gin.Context) {
		userID, ok := domain.AuditUserID(c.Request.Context())
		if !ok || userID != "user-id" {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("got status %d", recorder.Code)
	}
}
