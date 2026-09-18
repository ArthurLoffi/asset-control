package ports

import (
	"context"

	"asset-control/src/main/domain/dto"
)

type ChangePorts interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
	CreateChange(ctx context.Context, input dto.ChangeDto) error
}
