package entities

import "time"

type Machine struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	User             string    `gorm:"column:user;type:varchar(255);not null" json:"user"`
	AssetNumber      *int      `gorm:"column:asset_number" json:"asset_number"`
	Type             string    `gorm:"column:type;type:enum('Notebook','PC');not null" json:"type"`
	Model            string    `gorm:"column:model;type:varchar(255);not null" json:"model"`
	SerialNumber     string    `gorm:"column:serial_number;type:varchar(100);not null" json:"serial_number"`
	HasKeyboard      *bool     `gorm:"column:has_keyboard" json:"has_keyboard"`
	HasMouse         *bool     `gorm:"column:has_mouse" json:"has_mouse"`
	MachineCondition *string   `gorm:"column:machine_condition;type:enum('Novo','Bom','Ruim')" json:"machine_condition"`
	Situation        string    `gorm:"column:situation;type:enum('Em uso','Sem uso','Com defeito','Uso temporário','Manutenção');not null;default:'Em uso'" json:"situation"`
	EmployeeID       *uint64   `gorm:"column:employee_id" json:"employee_id"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Machine) TableName() string {
	return "machines"
}
