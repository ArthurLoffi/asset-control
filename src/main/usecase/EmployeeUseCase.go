package usecase

import (
	"context"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type EmployeeUseCase struct{ repository ports.EmployeesPorts }

func NewEmployeeUseCase(repository ports.EmployeesPorts) *EmployeeUseCase {
	return &EmployeeUseCase{repository: repository}
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
	return uc.repository.UpdateEmployee(ctx, id, input)
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
