package dto

import "time"

type MonitorDto struct {
	ID               uint64    `json:"id"`
	HasMonitor       *bool     `json:"has_monitor"`
	AssetNumber      *uint64   `json:"asset_number"`
	Brand            *string   `json:"brand"`
	Model            string    `json:"model"`
	SerialNumber     *string   `json:"serial_number"`
	MonitorCondition *string   `json:"monitor_condition"`
	Situation        string    `json:"situation"`
	EmployeeID       *uint64   `json:"employee_id"`
	EmployeeIDs      []uint64  `json:"employee_ids"`
	CreatedAt        time.Time `json:"created_at"`
}
