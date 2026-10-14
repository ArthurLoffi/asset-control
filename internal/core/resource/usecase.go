package resource

import (
	"context"
	"log/slog"

	"asset-control/internal/core/audit"
	"asset-control/internal/core/dto"
)

type FindAllUseCase[TRequest any, TResponse any] struct {
	port     ReadPort[TRequest, TResponse]
	resource string
	logger   *slog.Logger
}

func (u *FindAllUseCase[TRequest, TResponse]) Execute(ctx context.Context, request TRequest) (dto.PaginatedResponseDTO[TResponse], error) {
	u.logger.InfoContext(ctx, "[execute] Find all "+u.resource, "request", dto.Describe(request))
	return u.port.FindAll(ctx, request)
}

type FindByIDUseCase[TRequest any, TResponse any] struct {
	port     ReadPort[TRequest, TResponse]
	resource string
	logger   *slog.Logger
}

func (u *FindByIDUseCase[TRequest, TResponse]) Execute(ctx context.Context, id uint64) (TResponse, error) {
	u.logger.InfoContext(ctx, "[execute] Find "+u.resource+" by id", "id", id)
	return u.port.FindByID(ctx, id)
}

type CreateUseCase[TRequest any, TInput Validated[TInput], TResponse any] struct {
	port     DataPort[TRequest, TInput, TResponse]
	resource string
	logger   *slog.Logger
}

func (u *CreateUseCase[TRequest, TInput, TResponse]) Execute(ctx context.Context, input TInput) (TResponse, error) {
	u.logger.InfoContext(ctx, "[execute] Create "+u.resource)
	input, err := input.Validated()
	if err != nil {
		var zero TResponse
		return zero, err
	}
	return u.port.Create(ctx, input)
}

// UpdateUseCase records before/after snapshots in the same transaction as the write.
type UpdateUseCase[TRequest any, TInput Validated[TInput], TResponse any] struct {
	port     DataPort[TRequest, TInput, TResponse]
	recorder audit.Recorder
	resource string
	logger   *slog.Logger
}

func (u *UpdateUseCase[TRequest, TInput, TResponse]) Execute(ctx context.Context, id uint64, input TInput) (TResponse, error) {
	u.logger.InfoContext(ctx, "[execute] Update "+u.resource, "id", id)
	input, err := input.Validated()
	if err != nil {
		var zero TResponse
		return zero, err
	}
	return audit.UpdateWithChange(ctx, u.recorder, u.resource, int64(id),
		func(txCtx context.Context) (TResponse, error) { return u.port.FindByID(txCtx, id) },
		func(txCtx context.Context) (TResponse, error) { return u.port.Update(txCtx, id, input) },
	)
}

type DeactivateUseCase[TRequest any, TInput any, TResponse any] struct {
	port     DataPort[TRequest, TInput, TResponse]
	resource string
	logger   *slog.Logger
}

func (u *DeactivateUseCase[TRequest, TInput, TResponse]) Execute(ctx context.Context, id uint64) (dto.DeactivationResponseDTO, error) {
	u.logger.InfoContext(ctx, "[execute] Deactivate "+u.resource, "id", id)
	return u.port.Deactivate(ctx, id)
}

type FindInactiveUseCase[TRequest any, TInput any, TResponse any] struct {
	port     DataPort[TRequest, TInput, TResponse]
	resource string
	logger   *slog.Logger
}

func (u *FindInactiveUseCase[TRequest, TInput, TResponse]) Execute(ctx context.Context, request dto.PaginatedRequestDTO) (dto.PaginatedResponseDTO[dto.DeactivationResponseDTO], error) {
	u.logger.InfoContext(ctx, "[execute] Find inactive "+u.resource, "request", dto.Describe(request))
	return u.port.FindInactive(ctx, request)
}

type ReadUseCases[TRequest any, TResponse any] struct {
	FindAll  *FindAllUseCase[TRequest, TResponse]
	FindByID *FindByIDUseCase[TRequest, TResponse]
}

func NewReadUseCases[TRequest any, TResponse any](resource string, port ReadPort[TRequest, TResponse]) ReadUseCases[TRequest, TResponse] {
	logger := slog.Default().With("component", resource)
	return ReadUseCases[TRequest, TResponse]{
		FindAll:  &FindAllUseCase[TRequest, TResponse]{port: port, resource: resource, logger: logger},
		FindByID: &FindByIDUseCase[TRequest, TResponse]{port: port, resource: resource, logger: logger},
	}
}

type UseCases[TRequest any, TInput Validated[TInput], TResponse any] struct {
	ReadUseCases[TRequest, TResponse]
	Create       *CreateUseCase[TRequest, TInput, TResponse]
	Update       *UpdateUseCase[TRequest, TInput, TResponse]
	Deactivate   *DeactivateUseCase[TRequest, TInput, TResponse]
	FindInactive *FindInactiveUseCase[TRequest, TInput, TResponse]
}

// NewUseCases wires every action of one resource. resource is the singular name
// stored in audit records (for example "machine"). recorder may be nil in tests.
func NewUseCases[TRequest any, TInput Validated[TInput], TResponse any](resource string, port DataPort[TRequest, TInput, TResponse], recorder audit.Recorder) UseCases[TRequest, TInput, TResponse] {
	logger := slog.Default().With("component", resource)
	return UseCases[TRequest, TInput, TResponse]{
		ReadUseCases: NewReadUseCases(resource, port),
		Create:       &CreateUseCase[TRequest, TInput, TResponse]{port: port, resource: resource, logger: logger},
		Update:       &UpdateUseCase[TRequest, TInput, TResponse]{port: port, recorder: recorder, resource: resource, logger: logger},
		Deactivate:   &DeactivateUseCase[TRequest, TInput, TResponse]{port: port, resource: resource, logger: logger},
		FindInactive: &FindInactiveUseCase[TRequest, TInput, TResponse]{port: port, resource: resource, logger: logger},
	}
}
