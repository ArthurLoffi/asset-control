package usecase

import (
	"asset-control/internal/modules/organization/application/port"
	"asset-control/internal/modules/organization/domain/dto"
	"context"
)

type UtilsUseCase struct {
	utilsPorts port.UtilsPort
}

func NewUtilsUseCase(utilsPorts port.UtilsPort) *UtilsUseCase {
	return &UtilsUseCase{
		utilsPorts: utilsPorts,
	}
}

func (uuc *UtilsUseCase) GetAllBranchs(ctx context.Context) ([]dto.BranchDto, error) {
	return uuc.utilsPorts.GetAllBranchs(ctx)
}

func (uuc *UtilsUseCase) BranchIdToDescription(ctx context.Context, id uint64) (string, error) {
	return uuc.utilsPorts.BranchIdToDescription(ctx, id)
}

func (uuc *UtilsUseCase) GetAllDepartments(ctx context.Context) ([]dto.DepartmentDto, error) {
	return uuc.utilsPorts.GetAllDepartments(ctx)
}

func (uuc *UtilsUseCase) DepartmentIdToDescription(ctx context.Context, id uint64) (string, error) {
	return uuc.utilsPorts.DepartmentIdToDescription(ctx, id)
}
