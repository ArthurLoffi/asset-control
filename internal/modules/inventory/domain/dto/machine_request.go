package dto

import coredto "asset-control/internal/core/dto"

// MachineRequestDTO filters GET /machines.
type MachineRequestDTO struct {
	coredto.PaginatedRequestDTO
	AssignmentFilterDTO
}

var MachineSortableFields = []string{"id", "user", "model", "created_at"}

type MachineInputDTO struct {
	User             string   `json:"user" binding:"required,max=255"`
	AssetNumber      *int     `json:"asset_number" binding:"omitempty,gte=0"`
	Type             string   `json:"type" binding:"required,max=100"`
	Model            string   `json:"model" binding:"required,max=255"`
	SerialNumber     string   `json:"serial_number" binding:"max=100"`
	HasKeyboard      *bool    `json:"has_keyboard"`
	HasMouse         *bool    `json:"has_mouse"`
	MachineCondition *string  `json:"machine_condition" binding:"omitempty,max=100"`
	Situation        string   `json:"situation" binding:"omitempty,max=100"`
	EmployeeID       *uint64  `json:"employee_id" binding:"omitempty,gt=0,lte=9223372036854775807"`
	EmployeeIDs      []uint64 `json:"employee_ids" binding:"omitempty,dive,gt=0,lte=9223372036854775807"`
}

func (input MachineInputDTO) Validated() (MachineInputDTO, error) {
	err := input.Validate()
	return input, err
}
