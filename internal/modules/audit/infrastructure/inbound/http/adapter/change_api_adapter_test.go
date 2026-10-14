package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"asset-control/internal/modules/audit/application/usecase"
	"asset-control/internal/modules/audit/domain/dto"
	"github.com/gin-gonic/gin"
)

type recentChangeStub struct {
	changes []dto.RecentChangeDto
	err     error
	limit   int64
}

func (stub *recentChangeStub) ListRecent(_ context.Context, limit int64) ([]dto.RecentChangeDto, error) {
	stub.limit = limit
	return stub.changes, stub.err
}

func TestRecentChangeControllerReturnsBoundedSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &recentChangeStub{changes: []dto.RecentChangeDto{{
		ID: "abc", Time: time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC),
		UserID: "editor-id", UserName: "Ana", Resource: "machines", ResourceID: 42, Action: "update",
		ChangedFields: []dto.ChangedFieldDto{{Field: "model", Before: "Old", After: "New"}},
	}}}
	router := gin.New()
	router.GET("/api/changes/recent", NewRecentChangeController(usecase.NewRecentChangeUseCase(stub)).ListRecent)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/changes/recent", nil))
	if recorder.Code != http.StatusOK || stub.limit != 10 {
		t.Fatalf("status %d, limit %d", recorder.Code, stub.limit)
	}
	var response struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Data) != 1 || response.Data[0]["resource_id"] != float64(42) {
		t.Fatalf("unexpected response: %#v", response.Data)
	}
	if response.Data[0]["user_name"] != "Ana" {
		t.Fatalf("editor name missing: %#v", response.Data[0])
	}
	if _, exists := response.Data[0]["object"]; exists {
		t.Fatal("audit snapshots leaked into the dashboard response")
	}
}

func TestRecentChangeControllerHandlesDatabaseFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/changes/recent", NewRecentChangeController(usecase.NewRecentChangeUseCase(
		&recentChangeStub{err: errors.New("database unavailable")},
	)).ListRecent)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/changes/recent", nil))
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
}
