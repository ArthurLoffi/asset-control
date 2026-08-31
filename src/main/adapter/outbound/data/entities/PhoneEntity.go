package entities

import "time"

type Phone struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Equipament   *string   `gorm:"column:equipament;type:enum('Celular','Só chip','Sem chip')" json:"equipament"`
	Model        string    `gorm:"column:model;type:varchar(100);not null" json:"model"`
	IMEI         string    `gorm:"column:imei;type:varchar(100);not null" json:"imei"`
	SerialNumber string    `gorm:"column:serial_number;type:varchar(100);not null" json:"serial_number"`
	HasChip      *bool     `gorm:"column:has_chip" json:"has_chip"`
	Situation    string    `gorm:"column:situation;type:enum('Em uso','Sem uso','Com defeito','Uso temporário','Manutenção');not null;default:'Em uso'" json:"situation"`
	EmployeeID   *uint64   `gorm:"column:employee_id" json:"employee_id"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Phone) TableName() string {
	return "phones"
}
