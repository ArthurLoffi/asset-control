package dto

type CreateEmployeeDto struct {
	Name         string `json:"name" binding:"required,max=255"`
	Phone        string `json:"phone" binding:"omitempty,max=20"`
	Email        string `json:"email" binding:"required,email,max=255"`
	BranchID     uint64 `json:"branch_id" binding:"required,gt=0,lte=9223372036854775807"`
	DepartmentID uint64 `json:"department_id" binding:"required,gt=0,lte=9223372036854775807"`
}
