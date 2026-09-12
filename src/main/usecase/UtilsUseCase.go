package usecase

import "context"

import "asset-control/src/main/domain/ports"

type UtilsUseCase struct {
	utilsPorts ports.UtilsPort
}

func NewUtilsUseCase(utilsPorts ports.UtilsPort) *UtilsUseCase {
	return &UtilsUseCase{
		utilsPorts: utilsPorts,
	}
}

func (uuc *UtilsUseCase) BranchIdToDescription(ctx context.Context, id uint64) (string, error) {
	return uuc.utilsPorts.BranchIdToDescription(ctx, id)
}

func (uuc *UtilsUseCase) DepartmentIdToDescription(ctx context.Context, id uint64) (string, error) {
	return uuc.utilsPorts.DepartmentIdToDescription(ctx, id)
}
