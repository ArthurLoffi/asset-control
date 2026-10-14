package ports

import (
	"context"

	"asset-control/src/main/domain/dto"
)

type RecentChangePort interface {
	ListRecent(ctx context.Context, limit int64) ([]dto.RecentChangeDto, error)
}
