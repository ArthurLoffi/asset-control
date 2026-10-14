package usecase

import (
	"asset-control/internal/core/audit"
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
)

type EmployeeUseCase struct {
	repository port.EmployeesPorts
	changes    audit.Recorder
}

func NewEmployeeUseCase(repository port.EmployeesPorts, changePorts ...audit.Recorder) *EmployeeUseCase {
	var changes audit.Recorder
	if len(changePorts) > 0 {
		changes = changePorts[0]
	}
	return &EmployeeUseCase{repository: repository, changes: changes}
}

func (uc *EmployeeUseCase) CreateEmployee(ctx context.Context, input dto.CreateEmployeeDto) (dto.EmployeeDto, error) {
	if err := input.Validate(); err != nil {
		return dto.EmployeeDto{}, err
	}
	return uc.repository.CreateEmployee(ctx, input)
}

func (uc *EmployeeUseCase) UpdateEmployee(ctx context.Context, id uint64, input dto.CreateEmployeeDto) (dto.EmployeeDto, error) {
	if err := input.Validate(); err != nil {
		return dto.EmployeeDto{}, err
	}
	return audit.UpdateWithChange(
		ctx,
		uc.changes,
		"employee",
		int64(id),
		func(txCtx context.Context) (dto.EmployeeDto, error) { return uc.repository.FindByID(txCtx, id) },
		func(txCtx context.Context) (dto.EmployeeDto, error) {
			return uc.repository.UpdateEmployee(txCtx, id, input)
		},
	)
}

func (uc *EmployeeUseCase) GetAllEmployees(ctx context.Context) ([]dto.EmployeeDto, error) {
	return uc.repository.GetAllEmployees(ctx)
}

func (uc *EmployeeUseCase) FindByID(ctx context.Context, id uint64) (dto.EmployeeDto, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *EmployeeUseCase) FindByName(ctx context.Context, name string) (dto.EmployeeDto, error) {
	return uc.repository.FindByName(ctx, name)
}
