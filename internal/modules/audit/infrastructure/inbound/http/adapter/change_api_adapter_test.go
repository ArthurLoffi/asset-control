package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	coredto "asset-control/internal/core/dto"
	"asset-control/internal/modules/audit/application/usecase"
	"asset-control/internal/modules/audit/domain/dto"

	"github.com/gin-gonic/gin"
)

type fakeChangeDataPort struct {
	page            coredto.PaginatedResponseDTO[dto.ChangeResponseDTO]
	err             error
	receivedRequest dto.ChangeRequestDTO
}

func (f *fakeChangeDataPort) FindAll(_ context.Context, request dto.ChangeRequestDTO) (coredto.PaginatedResponseDTO[dto.ChangeResponseDTO], error) {
	f.receivedRequest = request
	return f.page, f.err
}

func serve(dataPort *fakeChangeDataPort, target string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewChangeAPIAdapter(usecase.NewFindAllChangesUseCase(dataPort)).RegisterRoutes(router)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))
	return recorder
}

func TestChangeAPIAdapterReturnsSummaryPage(t *testing.T) {
	dataPort := &fakeChangeDataPort{page: coredto.NewPaginatedResponseDTO([]dto.ChangeResponseDTO{{
		ID: "abc", Time: time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC),
		UserID: "editor-id", UserName: "Ana", Resource: "machine", ResourceID: 42, Action: "update",
		ChangedFields: []dto.ChangedFieldDTO{{Field: "model", Before: "Old", After: "New"}},
	}}, 0, 10, 1)}
	recorder := serve(dataPort, "/changes?page_size=10&direction=desc")
	if recorder.Code != http.StatusOK || dataPort.receivedRequest.PageSize != 10 || dataPort.receivedRequest.Direction != coredto.DirectionDesc {
		t.Fatalf("status %d, request %+v", recorder.Code, dataPort.receivedRequest)
	}
	var response struct {
		Content []map[string]any  `json:"content"`
		Paging  coredto.PagingDTO `json:"paging"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Content) != 1 || response.Content[0]["user_name"] != "Ana" || response.Paging.TotalElements != 1 {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
	if _, exists := response.Content[0]["object"]; exists {
		t.Fatal("audit snapshots leaked into the response")
	}
}

func TestChangeAPIAdapterRejectsUnknownSortAndHandlesFailures(t *testing.T) {
	if recorder := serve(&fakeChangeDataPort{}, "/changes?sort_by=user_id"); recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
	if recorder := serve(&fakeChangeDataPort{err: errors.New("database unavailable")}, "/changes"); recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
}
