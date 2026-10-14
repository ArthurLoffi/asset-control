package usecase

import (
	"context"

	"asset-control/internal/modules/audit/application/port"
	"asset-control/internal/modules/audit/domain/dto"
)

type RecentChangeUseCase struct{ repository port.RecentChangePort }

func NewRecentChangeUseCase(repository port.RecentChangePort) *RecentChangeUseCase {
	return &RecentChangeUseCase{repository: repository}
}

func (uc *RecentChangeUseCase) ListRecent(ctx context.Context) ([]dto.RecentChangeDto, error) {
	return uc.repository.ListRecent(ctx, 10)
}
