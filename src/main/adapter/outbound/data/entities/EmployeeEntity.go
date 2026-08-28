package entities

import "time"

type Employee struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Phone        string    `gorm:"column:phone;type:varchar(20);not null" json:"phone"`
	Email        string    `gorm:"column:email;type:varchar(255);not null" json:"email"`
	BranchID     uint64    `gorm:"column:branch_id;not null;" json:"branch_id"`
	DepartmentID uint64    `gorm:"column:department_id;not null;" json:"department_id"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Employee) TableName() string {
	return "employees"
}
