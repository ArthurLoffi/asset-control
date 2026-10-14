package resource

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"asset-control/internal/core/audit"
	"asset-control/internal/core/dto"
	"asset-control/internal/core/exception"
)

type fakeRequestDTO struct {
	dto.PaginatedRequestDTO
	Name *string `form:"name"`
}

type fakeInputDTO struct{ Name string }

func (input fakeInputDTO) Validated() (fakeInputDTO, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return input, exception.NewValidationError("name is required")
	}
	return input, nil
}

type fakeResponseDTO struct {
	ID   uint64
	Name string
}

// fakeDataPort implements DataPort without a database.
type fakeDataPort struct {
	receivedRequest fakeRequestDTO
	receivedInput   fakeInputDTO
	stored          fakeResponseDTO
}

func (f *fakeDataPort) FindAll(_ context.Context, request fakeRequestDTO) (dto.PaginatedResponseDTO[fakeResponseDTO], error) {
	f.receivedRequest = request
	return dto.NewPaginatedResponseDTO([]fakeResponseDTO{f.stored}, request.Page, 20, 1), nil
}

func (f *fakeDataPort) FindByID(_ context.Context, id uint64) (fakeResponseDTO, error) {
	if id != f.stored.ID {
		return fakeResponseDTO{}, exception.NotFoundByID("Fake", "FAKE_001", id)
	}
	return f.stored, nil
}

func (f *fakeDataPort) Create(_ context.Context, input fakeInputDTO) (fakeResponseDTO, error) {
	f.receivedInput = input
	return fakeResponseDTO{ID: 2, Name: input.Name}, nil
}

func (f *fakeDataPort) Update(_ context.Context, id uint64, input fakeInputDTO) (fakeResponseDTO, error) {
	f.receivedInput = input
	f.stored = fakeResponseDTO{ID: id, Name: input.Name}
	return f.stored, nil
}

func (f *fakeDataPort) Deactivate(_ context.Context, id uint64) (dto.DeactivationResponseDTO, error) {
	return dto.DeactivationResponseDTO{ID: id, DeactivatedAt: time.Unix(1, 0)}, nil
}

func (f *fakeDataPort) FindInactive(_ context.Context, request dto.PaginatedRequestDTO) (dto.PaginatedResponseDTO[dto.DeactivationResponseDTO], error) {
	return dto.NewPaginatedResponseDTO[dto.DeactivationResponseDTO](nil, request.Page, 20, 0), nil
}

// fakeRecorder runs the transaction inline and keeps the recorded change.
type fakeRecorder struct{ changes []audit.Change }

func (f *fakeRecorder) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func (f *fakeRecorder) CreateChange(_ context.Context, change audit.Change) error {
	f.changes = append(f.changes, change)
	return nil
}

func TestFindAllAndFindByIDDelegateToPort(t *testing.T) {
	dataPort := &fakeDataPort{stored: fakeResponseDTO{ID: 1, Name: "Ana"}}
	useCases := NewUseCases("fake", dataPort, nil)
	name := "an"
	request := fakeRequestDTO{PaginatedRequestDTO: dto.PaginatedRequestDTO{Page: 1}, Name: &name}

	page, err := useCases.FindAll.Execute(context.Background(), request)
	if err != nil || dataPort.receivedRequest != request || len(page.Content) != 1 {
		t.Fatalf("page %+v err %v", page, err)
	}
	if _, err := useCases.FindByID.Execute(context.Background(), 9); err == nil || !strings.Contains(err.Error(), "FAKE_001") {
		t.Fatalf("missing id must return the entity error, got %v", err)
	}
}

func TestCreateValidatesBeforePersisting(t *testing.T) {
	dataPort := &fakeDataPort{}
	useCases := NewUseCases("fake", dataPort, nil)
	if _, err := useCases.Create.Execute(context.Background(), fakeInputDTO{Name: " "}); !errors.As(err, new(*exception.ValidationError)) {
		t.Fatalf("expected validation error, got %v", err)
	}
	created, err := useCases.Create.Execute(context.Background(), fakeInputDTO{Name: " Ana "})
	if err != nil || created.Name != "Ana" || dataPort.receivedInput.Name != "Ana" {
		t.Fatalf("normalized input not persisted: %+v %v", created, err)
	}
}

func TestUpdateRecordsBeforeAndAfter(t *testing.T) {
	dataPort := &fakeDataPort{stored: fakeResponseDTO{ID: 1, Name: "Old"}}
	recorder := &fakeRecorder{}
	useCases := NewUseCases("fake", dataPort, recorder)
	if _, err := useCases.Update.Execute(context.Background(), 1, fakeInputDTO{Name: "New"}); err == nil {
		t.Fatal("updates without an authenticated user must fail")
	}
	ctx := audit.WithUserID(context.Background(), "user-id")
	updated, err := useCases.Update.Execute(ctx, 1, fakeInputDTO{Name: "New"})
	if err != nil || updated.Name != "New" || len(recorder.changes) != 1 {
		t.Fatalf("updated %+v err %v changes %d", updated, err, len(recorder.changes))
	}
	change := recorder.changes[0]
	if change.UserID != "user-id" || change.Resource != "fake" || change.ResourceID != 1 ||
		change.Object.Before.(fakeResponseDTO).Name != "Old" || change.Object.After.(fakeResponseDTO).Name != "New" {
		t.Fatalf("unexpected change %+v", change)
	}
}

func TestDeactivateAndFindInactiveDelegateToPort(t *testing.T) {
	useCases := NewUseCases("fake", &fakeDataPort{}, nil)
	if deactivation, err := useCases.Deactivate.Execute(context.Background(), 3); err != nil || deactivation.ID != 3 {
		t.Fatalf("%+v %v", deactivation, err)
	}
	if page, err := useCases.FindInactive.Execute(context.Background(), dto.PaginatedRequestDTO{}); err != nil || page.Content == nil {
		t.Fatalf("%+v %v", page, err)
	}
}
