package dto

type CreateMonitorDto struct {
	HasMonitor       *bool   `json:"has_monitor"`
	AssetNumber      *uint64 `json:"asset_number"`
	Brand            *string `json:"brand"`
	Model            string  `json:"model" binding:"required"`
	SerialNumber     *string `json:"serial_number"`
	MonitorCondition *string `json:"monitor_condition"`
	Situation        string  `json:"situation"`
	EmployeeID       *uint64 `json:"employee_id"`
}
