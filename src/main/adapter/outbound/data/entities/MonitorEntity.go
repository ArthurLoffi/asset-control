package entities

import "time"

type Monitor struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	HasMonitor       *bool     `gorm:"column:has_monitor" json:"has_monitor"`
	AssetNumber      *uint64   `gorm:"column:asset_number" json:"asset_number"`
	Brand            *string   `gorm:"column:brand;type:varchar(100)" json:"brand"`
	Model            string    `gorm:"column:model;type:varchar(100);not null" json:"model"`
	SerialNumber     *string   `gorm:"column:serial_number;type:varchar(100)" json:"serial_number"`
	MonitorCondition *string   `gorm:"column:monitor_condition;type:enum('Novo','Bom','Ruim')" json:"monitor_condition"`
	Situation        string    `gorm:"column:situation;type:enum('Em uso','Sem uso','Com defeito','Uso temporário','Manutenção');not null;default:'Em uso'" json:"situation"`
	EmployeeID       *uint64   `gorm:"column:employee_id" json:"employee_id"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Monitor) TableName() string {
	return "monitors"
}
