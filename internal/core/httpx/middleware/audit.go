package middleware

import (
	"asset-control/internal/core/audit"

	auth "github.com/ArthurLoffi/asset-control-auth"
	"github.com/gin-gonic/gin"
)

// AuditIdentity copies the already-validated authenticated user ID from Gin's
// session into the request context consumed by write use cases.
func AuditIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID, isAuthenticated := auth.UserID(c); isAuthenticated {
			c.Request = c.Request.WithContext(audit.WithUserID(c.Request.Context(), userID))
		}
		c.Next()
	}
}
