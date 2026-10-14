package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"asset-control/internal/core/exception"

	"github.com/gin-gonic/gin"
)

func TestWriteErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		err    error
		status int
		body   string
	}{
		{exception.New("Machine not found by id", "MACHINE_001", 404, map[string]any{"id": 7}), 404, `{"message":"Machine not found by id","code":"MACHINE_001","logs":{"id":7}}`},
		{fmt.Errorf("wrapped: %w", exception.NewValidationError("name is required")), 400, `{"message":"Request payload is invalid","code":"CORE_001","logs":{"validationErrors":["name is required"]}}`},
		{exception.ErrNotFound, 404, `{"message":"Record not found","code":"CORE_002","logs":{}}`},
		{exception.ErrConflict, 409, `{"message":"Record already exists","code":"CORE_004","logs":{}}`},
		{fmt.Errorf("query: %w", context.DeadlineExceeded), 504, `{"message":"Database operation timed out","code":"CORE_005","logs":{}}`},
		{errors.New("private connection details"), 500, `{"message":"Unexpected error","code":"CORE_003","logs":{}}`},
	} {
		router := gin.New()
		router.GET("/", func(c *gin.Context) { WriteError(c, tc.err) })
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
		if response.Code != tc.status || response.Body.String() != tc.body {
			t.Errorf("%v: %d %s", tc.err, response.Code, response.Body.String())
		}
	}
}

func TestNotFoundAndQueryID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(NotFound)
	router.GET("/item", func(c *gin.Context) {
		id, err := QueryID(c, "id")
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, id)
	})
	for target, want := range map[string]string{
		"/missing":                     `{"message":"Cannot GET /missing","code":"CORE_002","logs":{}}`,
		"/item":                        `{"message":"query param 'id' is required","code":"CORE_002","logs":{}}`,
		"/item?id=0":                   `{"message":"invalid 'id' value","code":"CORE_002","logs":{}}`,
		"/item?id=9223372036854775808": `{"message":"invalid 'id' value","code":"CORE_002","logs":{}}`,
		"/item?id=42":                  `42`,
	} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, target, nil))
		if response.Body.String() != want {
			t.Errorf("%s: %s", target, response.Body.String())
		}
	}
}
