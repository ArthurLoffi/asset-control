package dto

// AssignmentFilterDTO filters machines, monitors, and phones. Repeating
// employee_id matches records assigned to any of the employees; unassigned=true
// keeps only records without an employee.
type AssignmentFilterDTO struct {
	EmployeeIDs  []uint64 `form:"employee_id" binding:"dive,gt=0,lte=9223372036854775807"`
	IsUnassigned bool     `form:"unassigned"`
}
