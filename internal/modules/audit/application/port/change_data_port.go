package port

import (
	"context"

	coredto "asset-control/internal/core/dto"
	"asset-control/internal/modules/audit/domain/dto"
)

// ChangeDataPort is implemented by infrastructure/outbound/data/adapter.
type ChangeDataPort interface {
	FindAll(ctx context.Context, request dto.ChangeRequestDTO) (coredto.PaginatedResponseDTO[dto.ChangeResponseDTO], error)
}
