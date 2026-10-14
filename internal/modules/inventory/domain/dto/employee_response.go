package dto

import "time"

type EmployeeResponseDTO struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Phones       []string  `json:"phones"`
	Email        string    `json:"email"`
	Emails       []string  `json:"emails"`
	BranchID     uint64    `json:"branch_id"`
	DepartmentID uint64    `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
}
