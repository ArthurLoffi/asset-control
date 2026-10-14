package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"asset-control/internal/core/exception"

	auth "github.com/Brasil-Novo/asset-control-auth"
	"github.com/gin-gonic/gin"
)

// WriteError is the only way handlers report failures: every error leaves as
// {message, code, logs}. Unknown errors are opaque; details go to the log only.
func WriteError(c *gin.Context, err error) {
	cause := err
	if errors.Is(err, context.DeadlineExceeded) {
		err = exception.ErrTimeout
	}
	var domainError *exception.DomainError
	var validationError *exception.ValidationError
	var httpError *exception.HTTPError

	httpStatus := http.StatusInternalServerError
	body := exception.ErrorResponse{Message: "Unexpected error", Code: exception.UnexpectedCode, Logs: map[string]any{}}
	switch {
	case errors.As(err, &domainError):
		httpStatus = domainError.HTTPStatus
		body = exception.ErrorResponse{Message: domainError.Message, Code: domainError.Code, Logs: domainError.Logs}
	case errors.As(err, &validationError):
		httpStatus = http.StatusBadRequest
		body = exception.ErrorResponse{Message: "Request payload is invalid", Code: exception.ValidationCode,
			Logs: map[string]any{"validationErrors": validationError.Messages}}
	case errors.As(err, &httpError):
		httpStatus = httpError.HTTPStatus
		body = exception.ErrorResponse{Message: httpError.Message, Code: exception.HTTPCode, Logs: map[string]any{}}
	}

	summary := fmt.Sprintf("[catch] %d %s %s", httpStatus, body.Code, body.Message)
	if httpStatus >= http.StatusInternalServerError {
		slog.ErrorContext(c.Request.Context(), summary, "request_id", auth.RequestIDOf(c), "error", cause)
	} else {
		slog.WarnContext(c.Request.Context(), summary, "request_id", auth.RequestIDOf(c))
	}
	c.AbortWithStatusJSON(httpStatus, body)
}

// NotFound keeps unknown routes inside the error envelope.
func NotFound(c *gin.Context) {
	WriteError(c, exception.NewHTTPError(http.StatusNotFound, "Cannot "+c.Request.Method+" "+c.Request.URL.Path))
}

// BindJSON only decodes; use cases normalize and validate the input once.
func BindJSON(c *gin.Context, target any) error {
	if err := json.NewDecoder(c.Request.Body).Decode(target); err != nil {
		return exception.NewHTTPError(http.StatusBadRequest, "Malformed JSON body")
	}
	return nil
}
