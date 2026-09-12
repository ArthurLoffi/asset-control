package ports

import (
	"asset-control/src/main/domain/dto"
	"context"
)

type UtilsPort interface {
	GetAllBranchs(ctx context.Context) ([]dto.BranchDto, error)

	BranchIdToDescription(ctx context.Context, id uint64) (string, error)

	GetAllDepartments(ctx context.Context) ([]dto.DepartmentDto, error)

	DepartmentIdToDescription(ctx context.Context, id uint64) (string, error)
}
