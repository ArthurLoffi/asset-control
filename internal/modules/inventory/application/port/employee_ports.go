package port

import (
	"context"

	"asset-control/internal/modules/inventory/domain/dto"
)

type EmployeesPorts interface {
	CreateEmployee(ctx context.Context, input dto.CreateEmployeeDto) (dto.EmployeeDto, error)
	UpdateEmployee(ctx context.Context, id uint64, input dto.CreateEmployeeDto) (dto.EmployeeDto, error)
	GetAllEmployees(ctx context.Context) ([]dto.EmployeeDto, error)
	FindByID(ctx context.Context, id uint64) (dto.EmployeeDto, error)
	FindByName(ctx context.Context, name string) (dto.EmployeeDto, error)
}
