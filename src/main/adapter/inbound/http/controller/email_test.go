package controller

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"asset-control/src/main/domain"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/usecase"

	"github.com/gin-gonic/gin"
)

type emailStub struct {
	listResult []dto.EmailDto
	result     dto.EmailDto
	err        error
	ctx        context.Context
	id         int32
	user       string
	input      dto.CreateEmailDto
}

func (s *emailStub) GetAllEmails(ctx context.Context) ([]dto.EmailDto, error) {
	s.ctx = ctx
	return s.listResult, s.err
}

func TestEmailListResponseAndErrors(t *testing.T) {
	for _, records := range [][]dto.EmailDto{{}, {{ID: 42, User: "Ana", Email: "ana@example.com", PasswordHash: " plain password "}}} {
		stub := &emailStub{listResult: records}
		ctl := NewEmailController(usecase.NewEmailUseCase(stub))
		ctx, cancel := context.WithCancel(context.Background())
		result := request(ctl.GetAllEmailsController(), "GET", "/", "", ctx)
		cancel()
		var response dto.EmailListResponseDto
		if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if result.Code != 200 || response.Status != "emails retrieved successfully" || response.Data == nil || len(response.Data) != len(records) || stub.ctx != ctx {
			t.Fatal("incorrect email list response or context")
		}
		if len(records) > 0 && response.Data[0] != records[0] {
			t.Fatal("email list changed stored fields")
		}
	}
	for _, tc := range []struct {
		err    error
		status int
	}{{errors.New("private details"), 500}, {context.DeadlineExceeded, 504}} {
		ctl := NewEmailController(usecase.NewEmailUseCase(&emailStub{err: tc.err}))
		result := request(ctl.GetAllEmailsController(), "GET", "/", "", context.Background())
		if result.Code != tc.status || strings.Contains(result.Body.String(), "private details") {
			t.Fatal("incorrect email list error response")
		}
	}
}

func (s *emailStub) CreateEmail(ctx context.Context, input dto.CreateEmailDto) (dto.EmailDto, error) {
	s.ctx, s.input = ctx, input
	return s.result, s.err
}

func (s *emailStub) UpdateEmail(ctx context.Context, id int32, input dto.CreateEmailDto) (dto.EmailDto, error) {
	s.ctx, s.id, s.input = ctx, id, input
	return s.result, s.err
}

func TestEmailWrites(t *testing.T) {
	for _, method := range []string{"POST", "PUT"} {
		t.Run(method, func(t *testing.T) {
			want := dto.EmailDto{ID: 42, User: "Ana", Email: "ana@example.com", PasswordHash: " plain<&>text ", CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
			stub := &emailStub{result: want}
			ctl := NewEmailController(usecase.NewEmailUseCase(stub))
			handler, status, message := ctl.CreateEmailController(), 201, "created"
			if method == "PUT" {
				handler, status, message = ctl.UpdateEmailController(), 200, "email updated successfully"
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			body := `{"user":" Ana ","email":"ana@example.com","password_hash":" plain<&>text ","id":999,"created_at":"2000-01-01T00:00:00Z"}`
			result := request(handler, method, "/?id=42", body, ctx)
			var response dto.EmailResponseDto
			if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if result.Code != status || response.Status != message || response.Data != want {
				t.Fatal("write did not return the saved record and expected status")
			}
			if stub.ctx != ctx || stub.input.User != "Ana" || stub.input.PasswordHash != want.PasswordHash || (method == "PUT" && stub.id != 42) {
				t.Fatal("write input or context was not preserved correctly")
			}
			for _, tc := range []struct {
				err    error
				status int
			}{
				{domain.ErrConflict, 409}, {context.DeadlineExceeded, 504}, {errors.New("private details"), 500},
			} {
				stub.err = tc.err
				result := request(handler, method, "/?id=42", body, ctx)
				if result.Code != tc.status || strings.Contains(result.Body.String(), "private details") {
					t.Fatal("incorrect write error mapping")
				}
			}
			if method == "PUT" {
				stub.err = domain.ErrNotFound
				if result := request(handler, method, "/?id=42", body, ctx); result.Code != 404 {
					t.Fatal("missing email must return 404")
				}
			}
		})
	}
}

func TestInvalidEmailWritesDoNotReachRepository(t *testing.T) {
	stub := &emailStub{}
	ctl := NewEmailController(usecase.NewEmailUseCase(stub))
	for _, handler := range []gin.HandlerFunc{ctl.CreateEmailController(), ctl.UpdateEmailController()} {
		for _, body := range []string{
			`{broken`, `null`, `{}`, `{"user":"Ana","email":"ana@example.com"}`,
			`{"user":" ","email":"ana@example.com","password_hash":"secret"}`,
			`{"user":"Ana","email":"invalid","password_hash":"secret"}`,
			`{"user":"Ana","email":"ana@example.com","password_hash":""}`,
		} {
			if result := request(handler, "POST", "/?id=42", body, context.Background()); result.Code != 400 {
				t.Fatalf("invalid body accepted: status=%d", result.Code)
			}
		}
	}
	for _, id := range []string{"", "0", "-1", "abc", "2147483648"} {
		result := request(ctl.UpdateEmailController(), "PUT", "/?id="+id, `{"user":"Ana","email":"ana@example.com","password_hash":"secret"}`, context.Background())
		if result.Code != 400 {
			t.Fatalf("invalid update ID accepted: %q", id)
		}
	}
	if stub.ctx != nil {
		t.Fatal("invalid email write reached repository")
	}
}

func (s *emailStub) FindByID(ctx context.Context, id int32) (dto.EmailDto, error) {
	s.ctx, s.id = ctx, id
	return s.result, s.err
}

func (s *emailStub) FindByUser(ctx context.Context, user string) (dto.EmailDto, error) {
	s.ctx, s.user = ctx, user
	return s.result, s.err
}

func TestEmailInvalidSearchDoesNotQueryDatabase(t *testing.T) {
	stub := &emailStub{}
	ctl := NewEmailController(usecase.NewEmailUseCase(stub))
	for _, value := range []string{"", "0", "-1", "abc", "1.5", "2147483648", "9223372036854775808"} {
		result := request(ctl.FindByIDController(), "GET", "/?id="+value, "", context.Background())
		if result.Code != 400 {
			t.Fatalf("ID %q: %d", value, result.Code)
		}
	}
	for _, target := range []string{"/", "/?user=", "/?user=%20%09", "/?email=ana@example.com"} {
		result := request(ctl.FindByUserController(), "GET", target, "", context.Background())
		if result.Code != 400 {
			t.Fatalf("user query %q: %d", target, result.Code)
		}
	}
	if stub.ctx != nil {
		t.Fatal("invalid search reached repository")
	}
}

func TestEmailLookupResponseAndContext(t *testing.T) {
	want := dto.EmailDto{ID: 2147483647, User: "Ana", Email: "ana@example.com", PasswordHash: " plain<&>text password ", CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	stub := &emailStub{result: want}
	ctl := NewEmailController(usecase.NewEmailUseCase(stub))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for target, handler := range map[string]gin.HandlerFunc{
		"/?id=2147483647":  ctl.FindByIDController(),
		"/?user=%20Ana%20": ctl.FindByUserController(),
	} {
		result := request(handler, "GET", target, "", ctx)
		if result.Code != 200 {
			t.Fatalf("%s: status=%d", target, result.Code)
		}
		var response dto.EmailResponseDto
		if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Status != "email found successfully" || response.Data != want {
			t.Fatal("response did not preserve the stored email fields and plaintext password")
		}
		if stub.ctx != ctx {
			t.Fatal("request context was not propagated")
		}
	}
	if stub.id != want.ID || stub.user != "Ana" {
		t.Fatal("search values were not passed correctly")
	}
}

func TestEmailLookupErrors(t *testing.T) {
	for _, tc := range []struct {
		err    error
		status int
	}{
		{domain.ErrNotFound, 404},
		{errors.New("private database details"), 500},
		{context.DeadlineExceeded, 504},
	} {
		stub := &emailStub{err: tc.err}
		ctl := NewEmailController(usecase.NewEmailUseCase(stub))
		for target, handler := range map[string]gin.HandlerFunc{
			"/?id=1":     ctl.FindByIDController(),
			"/?user=Ana": ctl.FindByUserController(),
		} {
			result := request(handler, "GET", target, "", context.Background())
			if result.Code != tc.status || strings.Contains(result.Body.String(), "private database details") {
				t.Fatalf("%s: unexpected error response %d %s", target, result.Code, result.Body.String())
			}
		}
	}
}
