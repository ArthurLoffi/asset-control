package dto

type CreateMonitorDto struct {
	HasMonitor       *bool   `json:"has_monitor"`
	AssetNumber      *uint64 `json:"asset_number" binding:"omitempty,lte=9223372036854775807"`
	Brand            *string `json:"brand" binding:"omitempty,max=100"`
	Model            string  `json:"model" binding:"required,max=100"`
	SerialNumber     *string `json:"serial_number" binding:"omitempty,max=100"`
	MonitorCondition *string `json:"monitor_condition" binding:"omitempty,max=100"`
	Situation        string  `json:"situation" binding:"omitempty,max=100"`
	EmployeeID       *uint64 `json:"employee_id" binding:"omitempty,gt=0,lte=9223372036854775807"`
}
