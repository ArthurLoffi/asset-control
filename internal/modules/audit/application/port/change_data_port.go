package port

import (
	"context"

	"asset-control/internal/modules/audit/domain/dto"
)

type RecentChangePort interface {
	ListRecent(ctx context.Context, limit int64) ([]dto.RecentChangeDto, error)
}
