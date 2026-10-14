package schema

import "time"

const EmployeeCollection = "employees"

// EmployeeSortColumns maps the public sort_by value to the stored field.
var EmployeeSortColumns = map[string]string{"id": "id", "name": "name", "created_at": "created_at"}

type EmployeeDocument struct {
	ID           uint64    `bson:"id"`
	Name         string    `bson:"name"`
	Phone        string    `bson:"phone"`
	Phones       []string  `bson:"phones"`
	Email        string    `bson:"email"`
	Emails       []string  `bson:"emails"`
	BranchID     uint64    `bson:"branch_id"`
	DepartmentID uint64    `bson:"department_id"`
	CreatedAt    time.Time `bson:"created_at"`
}
