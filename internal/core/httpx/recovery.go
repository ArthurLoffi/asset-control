package httpx

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Recovery avoids Gin's request dump, which can include the session cookie, and
// never logs the panic value, which can include request data.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				WriteError(c, errors.New("unhandled panic"))
			}
		}()
		c.Next()
	}
}
