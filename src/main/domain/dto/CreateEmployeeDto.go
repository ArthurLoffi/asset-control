package dto

type CreateEmployeeDto struct {
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone"`
	Email        string `json:"email" binding:"required,email"`
	BranchID     uint64 `json:"branch_id" binding:"required"`
	DepartmentID uint64 `json:"department_id" binding:"required"`
}
