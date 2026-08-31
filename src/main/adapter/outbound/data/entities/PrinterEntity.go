package entities

import "time"

type Printer struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AssetNumber      uint64    `gorm:"column:asset_number;not null" json:"asset_number"`
	Model            string    `gorm:"column:model;type:varchar(100);not null" json:"model"`
	IP               string    `gorm:"column:ip;type:varchar(14);not null" json:"ip"`
	MachineCondition *string   `gorm:"column:machine_condition;type:enum('Novo','Bom','Ruim')" json:"machine_condition"`
	BranchID         uint64    `gorm:"column:branch_id;not null" json:"branch_id"`
	DepartmentsID    uint64    `gorm:"column:departments_id;not null" json:"departments_id"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Printer) TableName() string {
	return "printers"
}
