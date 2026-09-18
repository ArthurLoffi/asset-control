package usecase

import (
	"asset-control/src/main/domain"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
	"context"
	"math"
)

type DeactivationUseCase struct{ repository ports.DeactivationPort }

func NewDeactivationUseCase(repository ports.DeactivationPort) *DeactivationUseCase {
	return &DeactivationUseCase{repository: repository}
}

func validResourceID(resource string, id uint64) bool {
	switch resource {
	case "employees", "machines", "phones", "monitors", "printers":
		return id <= math.MaxInt64
	case "emails":
		return id <= math.MaxInt32
	default:
		return false
	}
}

func (uc *DeactivationUseCase) Deactivate(ctx context.Context, resource string, id uint64) (dto.DeactivationDto, error) {
	if id == 0 || !validResourceID(resource, id) {
		return dto.DeactivationDto{}, domain.ErrValidation
	}
	return uc.repository.Deactivate(ctx, resource, id)
}

func (uc *DeactivationUseCase) ListInactive(ctx context.Context, resource string, after uint64, limit int) (dto.InactivePageDto, error) {
	if !validResourceID(resource, after) || limit < 1 || limit > 100 {
		return dto.InactivePageDto{}, domain.ErrValidation
	}
	return uc.repository.ListInactive(ctx, resource, after, limit)
}
