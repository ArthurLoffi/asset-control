package dto

type CreatePhoneDto struct {
	Equipament   *string `json:"equipament" binding:"omitempty,max=100"`
	Model        string  `json:"model" binding:"required,max=100"`
	IMEI         string  `json:"imei" binding:"required,max=100"`
	SerialNumber string  `json:"serial_number" binding:"required,max=100"`
	HasChip      *bool   `json:"has_chip"`
	Situation    string  `json:"situation" binding:"omitempty,max=100"`
	EmployeeID   *uint64 `json:"employee_id" binding:"omitempty,gt=0,lte=9223372036854775807"`
}
