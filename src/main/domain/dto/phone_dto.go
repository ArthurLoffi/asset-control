package dto

import "time"

type PhoneDto struct {
	ID           uint64    `json:"id"`
	Equipament   *string   `json:"equipament"`
	Model        string    `json:"model"`
	IMEI         string    `json:"imei"`
	SerialNumber string    `json:"serial_number"`
	HasChip      *bool     `json:"has_chip"`
	Situation    string    `json:"situation"`
	EmployeeID   *uint64   `json:"employee_id"`
	CreatedAt    time.Time `json:"created_at"`
}
