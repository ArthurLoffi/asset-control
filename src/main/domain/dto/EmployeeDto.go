package dto

import "time"

type EmployeeDto struct {
	ID           uint
	Name         string
	Phone        string
	Email        string
	BranchID     uint
	DepartmentID uint
	CreatedAt    time.Time
}
