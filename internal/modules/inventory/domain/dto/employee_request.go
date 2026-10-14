package dto

import coredto "asset-control/internal/core/dto"

// EmployeeRequestDTO filters GET /employees. Name matches partially and
// case-insensitively; the other filters are exact.
type EmployeeRequestDTO struct {
	coredto.PaginatedRequestDTO
	Name         *string `form:"name"`
	BranchID     *uint64 `form:"branch_id"`
	DepartmentID *uint64 `form:"department_id"`
}

var EmployeeSortableFields = []string{"id", "name", "created_at"}

type EmployeeInputDTO struct {
	Name         string   `json:"name" binding:"required,max=255"`
	Phone        string   `json:"phone" binding:"omitempty,max=20"`
	Phones       []string `json:"phones" binding:"omitempty,dive,max=20"`
	Email        string   `json:"email" binding:"omitempty,email,max=255"`
	Emails       []string `json:"emails" binding:"omitempty,dive,email,max=255"`
	BranchID     uint64   `json:"branch_id" binding:"required,gt=0,lte=9223372036854775807"`
	DepartmentID uint64   `json:"department_id" binding:"required,gt=0,lte=9223372036854775807"`
}

func (input EmployeeInputDTO) Validated() (EmployeeInputDTO, error) {
	err := input.Validate()
	return input, err
}
