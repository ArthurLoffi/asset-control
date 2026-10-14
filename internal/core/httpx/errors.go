package httpx

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	exception "asset-control/internal/core/exception"

	"github.com/gin-gonic/gin"
)

// BSON integers are signed; reject values that the driver cannot encode.
func ParseID(value string) (uint64, error) {
	id, err := strconv.ParseUint(value, 10, 63)
	if err != nil || id == 0 {
		return 0, exception.ErrValidation
	}
	return id, nil
}

func WriteError(c *gin.Context, err error) {
	status, message := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, exception.ErrValidation):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, exception.ErrNotFound):
		status, message = http.StatusNotFound, exception.ErrNotFound.Error()
	case errors.Is(err, exception.ErrConflict):
		status, message = http.StatusConflict, exception.ErrConflict.Error()
	case errors.Is(err, context.DeadlineExceeded):
		status, message = http.StatusGatewayTimeout, "database operation timed out"
	}
	if status >= 500 {
		_ = c.Error(err)
	}
	c.JSON(status, gin.H{"error": message})
}
