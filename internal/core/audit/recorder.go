package audit

import (
	"context"
	"errors"
	"time"
)

var errAuditUserMissing = errors.New("authenticated user is missing from audit context")

func UpdateWithChange[T any](
	ctx context.Context,
	changes Recorder,
	resource string,
	resourceID int64,
	find func(context.Context) (T, error),
	update func(context.Context) (T, error),
) (T, error) {
	if changes == nil {
		return update(ctx)
	}

	userID, ok := UserID(ctx)
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

		return changes.CreateChange(txCtx, Change{
			Time:       time.Now().UTC(),
			UserID:     userID,
			Resource:   resource,
			ResourceID: resourceID,
			Action:     "update",
			Object: ChangeObject{
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
