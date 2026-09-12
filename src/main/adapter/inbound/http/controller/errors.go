package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"asset-control/src/main/domain"

	"github.com/gin-gonic/gin"
)

// BSON integers are signed; reject values that the driver cannot encode.
func parseID(value string) (uint64, error) {
	id, err := strconv.ParseUint(value, 10, 63)
	if err != nil || id == 0 {
		return 0, domain.ErrValidation
	}
	return id, nil
}

func writeError(c *gin.Context, err error) {
	status, message := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, domain.ErrValidation):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, message = http.StatusNotFound, domain.ErrNotFound.Error()
	case errors.Is(err, domain.ErrConflict):
		status, message = http.StatusConflict, domain.ErrConflict.Error()
	case errors.Is(err, context.DeadlineExceeded):
		status, message = http.StatusGatewayTimeout, "database operation timed out"
	}
	if status >= 500 {
		_ = c.Error(err)
	}
	c.JSON(status, gin.H{"error": message})
}
