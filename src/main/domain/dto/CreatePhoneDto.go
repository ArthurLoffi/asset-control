package dto

type CreatePhoneDto struct {
	Equipament   *string `json:"equipament"`
	Model        string  `json:"model" binding:"required"`
	IMEI         string  `json:"imei" binding:"required"`
	SerialNumber string  `json:"serial_number" binding:"required"`
	HasChip      *bool   `json:"has_chip"`
	Situation    string  `json:"situation"`
	EmployeeID   *uint64 `json:"employee_id"`
}
