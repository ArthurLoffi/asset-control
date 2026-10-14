package schema

import "time"

type Employee struct {
	ID           uint64    `bson:"id" json:"id"`
	Name         string    `bson:"name" json:"name"`
	Phone        string    `bson:"phone" json:"phone"`
	Phones       []string  `bson:"phones" json:"phones"`
	Email        string    `bson:"email" json:"email"`
	Emails       []string  `bson:"emails" json:"emails"`
	BranchID     uint64    `bson:"branch_id" json:"branch_id"`
	DepartmentID uint64    `bson:"department_id" json:"department_id"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

func (Employee) CollectionName() string {
	return "employees"
}
