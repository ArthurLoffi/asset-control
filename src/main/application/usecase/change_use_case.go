package usecase

import (
	"context"
	"errors"
	"time"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/domain"
	"asset-control/src/main/domain/dto"
)

var errAuditUserMissing = errors.New("authenticated user is missing from audit context")

func updateWithChange[T any](
	ctx context.Context,
	changes ports.ChangePorts,
	resource string,
	resourceID int64,
	find func(context.Context) (T, error),
	update func(context.Context) (T, error),
) (T, error) {
	if changes == nil {
		return update(ctx)
	}

	userID, ok := domain.AuditUserID(ctx)
	if !ok {
		var zero T
		return zero, errAuditUserMissing
	}

	var updated T
	err := changes.WithinTransaction(ctx, func(txCtx context.Context) error {
		before, err := find(txCtx)
		if err != nil {
			return err
		}

		after, err := update(txCtx)
		if err != nil {
			return err
		}
		updated = after

		return changes.CreateChange(txCtx, dto.ChangeDto{
			Time:       time.Now().UTC(),
			UserID:     userID,
			Resource:   resource,
			ResourceID: resourceID,
			Action:     "update",
			Object: dto.ChangeObjectDto{
				Before: before,
				After:  after,
			},
		})
	})
	if err != nil {
		var zero T
		return zero, err
	}

	return updated, nil
}
