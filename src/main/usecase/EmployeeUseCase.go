package usecase

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type EmployeeUseCase struct {
	employeePorts ports.EmployeesPorts
}

func NewEmployeeUseCase(employeePorts ports.EmployeesPorts) *EmployeeUseCase {
	return &EmployeeUseCase{
		employeePorts: employeePorts,
	}
}

func (uc *EmployeeUseCase) CreateEmployee(employee dto.EmployeeDto) error {
	return uc.employeePorts.CreateEmployee(employee)
}
