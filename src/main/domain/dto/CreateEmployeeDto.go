package dto

type CreateEmployeeDto struct {
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone"`
	Email        string `json:"email" binding:"required,email"`
	BranchID     uint   `json:"branch_id" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}
