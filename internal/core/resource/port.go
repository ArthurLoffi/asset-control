// Package resource holds the ports and use cases shared by every inventory-like
// entity: paginated listing, lookup by id, create, audited update, and soft delete.
package resource

import (
	"context"

	"asset-control/internal/core/dto"
)

// ReadPort is implemented by a driven adapter under infrastructure/outbound/data.
type ReadPort[TRequest any, TResponse any] interface {
	FindAll(ctx context.Context, request TRequest) (dto.PaginatedResponseDTO[TResponse], error)
	FindByID(ctx context.Context, id uint64) (TResponse, error)
}

// DataPort adds writes and the deactivation lifecycle to ReadPort.
type DataPort[TRequest any, TInput any, TResponse any] interface {
	ReadPort[TRequest, TResponse]
	Create(ctx context.Context, input TInput) (TResponse, error)
	Update(ctx context.Context, id uint64, input TInput) (TResponse, error)
	Deactivate(ctx context.Context, id uint64) (dto.DeactivationResponseDTO, error)
	FindInactive(ctx context.Context, request dto.PaginatedRequestDTO) (dto.PaginatedResponseDTO[dto.DeactivationResponseDTO], error)
}

// Validated normalizes the input (trims, defaults) and validates it, returning
// the value to persist.
type Validated[TInput any] interface {
	Validated() (TInput, error)
}
