package port

import (
	"asset-control/internal/modules/inventory/domain/dto"
	"context"
)

type DeactivationPort interface {
	Deactivate(context.Context, string, uint64) (dto.DeactivationDto, error)
	ListInactive(context.Context, string, uint64, int) (dto.InactivePageDto, error)
}
