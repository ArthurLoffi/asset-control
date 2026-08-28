package ports

import "asset-control/src/main/domain/dto"

type EmployeesPorts interface {
	CreateEmployee(employee dto.EmployeeDto) error

	FindByName(name string) (dto.EmployeeDto, error)
}
