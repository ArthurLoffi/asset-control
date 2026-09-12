package ports

import "context"

type UtilsPort interface {
	BranchIdToDescription(ctx context.Context, id uint64) (string, error)

	DepartmentIdToDescription(ctx context.Context, id uint64) (string, error)
}
