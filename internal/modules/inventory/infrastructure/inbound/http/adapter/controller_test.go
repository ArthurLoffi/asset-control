package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"asset-control/internal/core/exception"
	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/application/usecase"
	"asset-control/internal/modules/inventory/domain/dto"
	"github.com/gin-gonic/gin"
)

type employeeStub struct {
	port.EmployeesPorts
	err    error
	called bool
	ctx    context.Context
}

func (s *employeeStub) FindByID(ctx context.Context, id uint64) (dto.EmployeeDto, error) {
	s.called, s.ctx = true, ctx
	return dto.EmployeeDto{ID: id}, s.err
}

func (s *employeeStub) FindByName(ctx context.Context, name string) (dto.EmployeeDto, error) {
	s.called = true
	return dto.EmployeeDto{}, s.err
}

func (s *employeeStub) CreateEmployee(ctx context.Context, input dto.CreateEmployeeDto) (dto.EmployeeDto, error) {
	s.called, s.ctx = true, ctx
	return dto.EmployeeDto{ID: 42, Name: input.Name, CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, s.err
}

func request(handler gin.HandlerFunc, method, target, body string, ctx context.Context) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Any("/", handler)
	req := httptest.NewRequest(method, target, strings.NewReader(body)).WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	return recorder
}

func TestAllCreateEndpointsRejectInvalidRequests(t *testing.T) {
	handlers := map[string]gin.HandlerFunc{
		"employees": NewEmployeeController(usecase.NewEmployeeUseCase(nil)).CreateEmployeeController(),
		"machines":  NewMachineController(usecase.NewMachineUseCase(nil)).CreateMachineController(),
		"monitors":  NewMonitorController(usecase.NewMonitorUseCase(nil)).CreateMonitorController(),
		"phones":    NewPhoneController(usecase.NewPhoneUseCase(nil)).CreatePhoneController(),
		"printers":  NewPrinterController(usecase.NewPrinterUseCase(nil)).CreatePrinterController(),
	}
	for name, handler := range handlers {
		t.Run(name, func(t *testing.T) {
			for _, body := range []string{"{}", "null", "{broken"} {
				result := request(handler, http.MethodPost, "/", body, context.Background())
				if result.Code != http.StatusBadRequest {
					t.Fatalf("%s: status=%d body=%s", body, result.Code, result.Body.String())
				}
			}
		})
	}
}

func TestEmployeeErrorMapping(t *testing.T) {
	for _, tt := range []struct {
		name   string
		err    error
		status int
	}{
		{"missing", fmt.Errorf("wrapped: %w", exception.ErrNotFound), 404},
		{"database failure", errors.New("private connection details"), 500},
		{"timeout", context.DeadlineExceeded, 504},
		{"duplicate", exception.ErrConflict, 409},
		{"success", nil, 200},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub := &employeeStub{err: tt.err}
			ctl := NewEmployeeController(usecase.NewEmployeeUseCase(stub))
			result := request(ctl.FindByIDController(), "GET", "/?id=1", "", context.Background())
			if result.Code != tt.status {
				t.Fatalf("got %d: %s", result.Code, result.Body.String())
			}
			if strings.Contains(result.Body.String(), "private connection") {
				t.Fatal("database error leaked to client")
			}
		})
	}
}

func TestInvalidIDsAndMissingNameDoNotQueryDatabase(t *testing.T) {
	stub := &employeeStub{}
	ctl := NewEmployeeController(usecase.NewEmployeeUseCase(stub))
	for _, value := range []string{"", "0", "-1", "abc", "9223372036854775808", "18446744073709551615"} {
		result := request(ctl.FindByIDController(), "GET", "/?id="+value, "", context.Background())
		if result.Code != 400 {
			t.Fatalf("ID %q: %d", value, result.Code)
		}
	}
	for _, target := range []string{"/", "/?name=%20%20"} {
		result := request(ctl.FindByNameController(), "GET", target, "", context.Background())
		if result.Code != 400 {
			t.Fatalf("missing name: %d", result.Code)
		}
	}
	if stub.called {
		t.Fatal("invalid request reached repository")
	}
}

func TestCreateReturnsSavedRecordAndPropagatesContext(t *testing.T) {
	stub := &employeeStub{}
	ctl := NewEmployeeController(usecase.NewEmployeeUseCase(stub))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := request(ctl.CreateEmployeeController(), "POST", "/", `{"name":"Ana","email":"ana@example.com","branch_id":1,"department_id":2,"id":999}`, ctx)
	if result.Code != 201 {
		t.Fatalf("got %d: %s", result.Code, result.Body.String())
	}
	var response struct {
		Data dto.EmployeeDto `json:"data"`
	}
	if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Data.ID != 42 || response.Data.CreatedAt.IsZero() {
		t.Fatalf("not saved record: %+v", response.Data)
	}
	if stub.ctx != ctx {
		t.Fatal("request context was not propagated")
	}
}
