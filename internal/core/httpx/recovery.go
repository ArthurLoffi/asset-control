package httpx

import (
	"log/slog"
	"net/http"

	auth "github.com/Brasil-Novo/asset-control-auth"
	"github.com/gin-gonic/gin"
)

// Recovery avoids Gin's request dump, which can include the session cookie.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				slog.ErrorContext(c.Request.Context(), "[catch] unhandled request error", "request_id", auth.RequestIDOf(c))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
