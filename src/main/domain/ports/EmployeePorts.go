package ports

import "asset-control/src/main/domain/dto"

type EmployeesPorts interface {
	CreateEmployee(employee dto.EmployeeDto) error

	UpdateEmployee(id uint64, employee dto.EmployeeDto) (dto.EmployeeDto, error)

	GetAllEmployees() ([]dto.EmployeeDto, error)

	FindByName(name string) (dto.EmployeeDto, error)

	FindByID(id uint64) (dto.EmployeeDto, error)
}
