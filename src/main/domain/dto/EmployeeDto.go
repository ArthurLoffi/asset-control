package dto

import "time"

type EmployeeDto struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	BranchID     uint      `json:"branch_id"`
	DepartmentID uint      `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
}
