package repository

import (
	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"errors"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (er *EmployeeRepository) CreateEmployee(employee dto.EmployeeDto) error {
	employeeEntity := entities.Employee{
		Name:         employee.Name,
		Phone:        employee.Phone,
		Email:        employee.Email,
		BranchID:     employee.BranchID,
		DepartmentID: employee.DepartmentID,
	}

	if err := er.db.Create(&employeeEntity).Error; err != nil {
		return err
	}

	return nil
}

func (er *EmployeeRepository) UpdateEmployee(id uint64, employee dto.EmployeeDto) (dto.EmployeeDto, error) {
	var employeeEntity entities.Employee

	if err := er.db.First(&employeeEntity, id).Error; err != nil {
		return dto.EmployeeDto{}, err
	}

	employeeEntity.Name = employee.Name
	employeeEntity.Phone = employee.Phone
	employeeEntity.Email = employee.Email
	employeeEntity.BranchID = employee.BranchID
	employeeEntity.DepartmentID = employee.DepartmentID

	if err := er.db.Save(&employeeEntity).Error; err != nil {
		return dto.EmployeeDto{}, err
	}

	return dto.EmployeeDto{
		ID:           employeeEntity.ID,
		Name:         employeeEntity.Name,
		Phone:        employeeEntity.Phone,
		Email:        employeeEntity.Email,
		BranchID:     employeeEntity.BranchID,
		DepartmentID: employeeEntity.DepartmentID,
		CreatedAt:    employeeEntity.CreatedAt,
	}, nil
}

func (er *EmployeeRepository) GetAllEmployees() ([]dto.EmployeeDto, error) {
	var employeeEntities []entities.Employee

	if err := er.db.Find(&employeeEntities).Error; err != nil {
		return nil, err
	}

	employees := make([]dto.EmployeeDto, 0, len(employeeEntities))
	for _, e := range employeeEntities {
		employees = append(employees, dto.EmployeeDto{
			ID:           e.ID,
			Name:         e.Name,
			Phone:        e.Phone,
			Email:        e.Email,
			BranchID:     e.BranchID,
			DepartmentID: e.DepartmentID,
			CreatedAt:    e.CreatedAt,
		})
	}

	return employees, nil
}

func (er *EmployeeRepository) FindByID(id uint64) (dto.EmployeeDto, error) {
	var employeeEntity entities.Employee

	if err := er.db.First(&employeeEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.EmployeeDto{}, err
		}
		return dto.EmployeeDto{}, err
	}

	return dto.EmployeeDto{
		ID:           employeeEntity.ID,
		Name:         employeeEntity.Name,
		Phone:        employeeEntity.Phone,
		Email:        employeeEntity.Email,
		BranchID:     employeeEntity.BranchID,
		DepartmentID: employeeEntity.DepartmentID,
		CreatedAt:    employeeEntity.CreatedAt,
	}, nil
}

func (er *EmployeeRepository) FindByName(name string) (dto.EmployeeDto, error) {
	var employeeEntity entities.Employee

	if err := er.db.Where("name = ?", name).First(&employeeEntity).Error; err != nil {
		return dto.EmployeeDto{}, err
	}

	return dto.EmployeeDto{
		ID:           employeeEntity.ID,
		Name:         employeeEntity.Name,
		Phone:        employeeEntity.Phone,
		Email:        employeeEntity.Email,
		BranchID:     employeeEntity.BranchID,
		DepartmentID: employeeEntity.DepartmentID,
		CreatedAt:    employeeEntity.CreatedAt,
	}, nil
}
