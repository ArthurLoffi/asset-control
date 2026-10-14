package usecase

import (
	"context"
	"log/slog"

	coredto "asset-control/internal/core/dto"
	"asset-control/internal/modules/audit/application/port"
	"asset-control/internal/modules/audit/domain/dto"
)

type FindAllChangesUseCase struct {
	changeDataPort port.ChangeDataPort
	logger         *slog.Logger
}

func NewFindAllChangesUseCase(changeDataPort port.ChangeDataPort) *FindAllChangesUseCase {
	return &FindAllChangesUseCase{changeDataPort: changeDataPort, logger: slog.Default().With("component", "FindAllChangesUseCase")}
}

func (u *FindAllChangesUseCase) Execute(ctx context.Context, request dto.ChangeRequestDTO) (coredto.PaginatedResponseDTO[dto.ChangeResponseDTO], error) {
	u.logger.InfoContext(ctx, "[execute] Find all changes", "request", coredto.Describe(request))
	return u.changeDataPort.FindAll(ctx, request)
}
