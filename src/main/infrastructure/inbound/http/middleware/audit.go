package middleware

import (
	"asset-control/src/main/domain"

	authMiddleware "github.com/Brasil-Novo/asset-control-auth/src/main/adapter/inbound/http/middleware"
	authDto "github.com/Brasil-Novo/asset-control-auth/src/main/domain/dto"
	"github.com/gin-gonic/gin"
)

// AuditIdentity copies the already-validated authenticated user ID from Gin's
// session into the request context consumed by write use cases.
func AuditIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(authMiddleware.SessionKey)
		if exists {
			if session, ok := value.(authDto.SessionDto); ok && session.ID != "" {
				c.Request = c.Request.WithContext(domain.WithAuditUserID(c.Request.Context(), session.ID))
			}
		}
		c.Next()
	}
}
