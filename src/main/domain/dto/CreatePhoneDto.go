package dto

type CreatePhoneDto struct {
	Equipament *string `json:"equipament" binding:"omitempty,max=100"`
	// Required unless equipament is "Só chip" or "Sem chip".
	Model string `json:"model" binding:"max=100"`
	// Required unless equipament is "Só chip" or "Sem chip".
	IMEI string `json:"imei" binding:"max=100"`
	// Required unless equipament is "Só chip" or "Sem chip".
	SerialNumber string  `json:"serial_number" binding:"max=100"`
	HasChip      *bool   `json:"has_chip"`
	Situation    string  `json:"situation" binding:"omitempty,max=100"`
	EmployeeID   *uint64 `json:"employee_id" binding:"omitempty,gt=0,lte=9223372036854775807"`
}
