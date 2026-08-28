package repository

import (
	"asset-control/src/main/domain/dto"
	"database/sql"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (er *EmployeeRepository) CreateEmployee(employee dto.EmployeeDto) error {

	return nil
}

func (er *EmployeeRepository) FindByName(name string) (dto.EmployeeDto, error) {
	var employee dto.EmployeeDto

	return employee, nil
}
