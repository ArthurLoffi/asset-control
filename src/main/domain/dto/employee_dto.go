package dto

import "time"

type EmployeeDto struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	BranchID     uint64    `json:"branch_id"`
	DepartmentID uint64    `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
}
