package ports

import "asset-control/src/main/domain/dto"

type EmployeesPorts interface {
	createEmployee(employee dto.EmployeeDto) error

	findByName(name string) (dto.EmployeeDto, error)
}
