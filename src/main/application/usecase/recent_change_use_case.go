package usecase

import (
	"context"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/domain/dto"
)

type RecentChangeUseCase struct{ repository ports.RecentChangePort }

func NewRecentChangeUseCase(repository ports.RecentChangePort) *RecentChangeUseCase {
	return &RecentChangeUseCase{repository: repository}
}

func (uc *RecentChangeUseCase) ListRecent(ctx context.Context) ([]dto.RecentChangeDto, error) {
	return uc.repository.ListRecent(ctx, 10)
}
