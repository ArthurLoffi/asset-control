package usecase

import "asset-control/src/main/domain/ports"

type UtilsUseCase struct {
	utilsPorts ports.UtilsPort
}

func NewUtilsUseCase(utilsPorts ports.UtilsPort) *UtilsUseCase {
	return &UtilsUseCase{
		utilsPorts: utilsPorts,
	}
}

func (uuc *UtilsUseCase) BranchIdToDescription(id uint64) (string, error) {
	return uuc.utilsPorts.BranchIdToDescription(id)
}

func (uuc *UtilsUseCase) DepartmentIdToDescription(id uint64) (string, error) {
	return uuc.utilsPorts.DepartmentIdToDescription(id)
}
