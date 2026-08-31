package dto

type CreateMachineDto struct {
	User             string  `json:"user" binding:"required"`
	AssetNumber      *int    `json:"asset_number"`
	Type             string  `json:"type" binding:"required"`
	Model            string  `json:"model" binding:"required"`
	SerialNumber     string  `json:"serial_number" binding:"required"`
	HasKeyboard      *bool   `json:"has_keyboard"`
	HasMouse         *bool   `json:"has_mouse"`
	MachineCondition *string `json:"machine_condition"`
	Situation        string  `json:"situation"`
	EmployeeID       *uint64 `json:"employee_id"`
}
