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

func (uc *EmployeeUseCase) UpdateEmployee(id uint64, employee dto.EmployeeDto) (dto.EmployeeDto, error) {
	return uc.employeePorts.UpdateEmployee(id, employee)
}

func (uc *EmployeeUseCase) GetAllEmployees() ([]dto.EmployeeDto, error) {
	return uc.employeePorts.GetAllEmployees()
}

func (uc *EmployeeUseCase) FindByID(id uint64) (dto.EmployeeDto, error) {
	return uc.employeePorts.FindByID(id)
}

func (uc *EmployeeUseCase) FindByName(name string) (dto.EmployeeDto, error) {
	return uc.employeePorts.FindByName(name)
}
