package dto

import "time"

type MachineDto struct {
	ID               uint64    `json:"id"`
	User             string    `json:"user"`
	AssetNumber      *int      `json:"asset_number"`
	Type             string    `json:"type"`
	Model            string    `json:"model"`
	SerialNumber     string    `json:"serial_number"`
	HasKeyboard      *bool     `json:"has_keyboard"`
	HasMouse         *bool     `json:"has_mouse"`
	MachineCondition *string   `json:"machine_condition"`
	Situation        string    `json:"situation"`
	EmployeeID       *uint64   `json:"employee_id"`
	CreatedAt        time.Time `json:"created_at"`
}
