package exception

import (
	"net/http"
	"strings"
)

// ErrorResponse is the single error body: {message, code, logs}. logs is public;
// never put secrets in it.
type ErrorResponse struct {
	Message string         `json:"message"`
	Code    string         `json:"code"`
	Logs    map[string]any `json:"logs"`
}

type DomainError struct {
	Message    string
	Code       string
	HTTPStatus int
	Logs       map[string]any
}

func New(message string, code string, httpStatus int, logs map[string]any) *DomainError {
	if logs == nil {
		logs = map[string]any{}
	}
	return &DomainError{Message: message, Code: code, HTTPStatus: httpStatus, Logs: logs}
}

func (e *DomainError) Error() string { return e.Code + ": " + e.Message }

// ValidationError becomes CORE_001 with logs.validationErrors.
type ValidationError struct{ Messages []string }

func NewValidationError(messages ...string) *ValidationError {
	return &ValidationError{Messages: messages}
}

func (e *ValidationError) Error() string {
	return "validation failed: " + strings.Join(e.Messages, "; ")
}

// HTTPError becomes CORE_002: unknown route, malformed JSON, invalid ID.
type HTTPError struct {
	HTTPStatus int
	Message    string
}

func NewHTTPError(httpStatus int, message string) *HTTPError {
	return &HTTPError{HTTPStatus: httpStatus, Message: message}
}

func (e *HTTPError) Error() string { return e.Message }

const (
	ValidationCode = "CORE_001"
	HTTPCode       = "CORE_002"
	UnexpectedCode = "CORE_003"
	ConflictCode   = "CORE_004"
	TimeoutCode    = "CORE_005"
)

var (
	ErrNotFound = NewHTTPError(http.StatusNotFound, "Record not found")
	ErrConflict = New("Record already exists", ConflictCode, http.StatusConflict, nil)
	ErrTimeout  = New("Database operation timed out", TimeoutCode, http.StatusGatewayTimeout, nil)
)

// NotFoundByID builds the 404 domain error of an entity, e.g. ("Machine", "MACHINE_001").
func NotFoundByID(entity string, code string, id uint64) error {
	return New(entity+" not found by id", code, http.StatusNotFound, map[string]any{"id": id})
}
