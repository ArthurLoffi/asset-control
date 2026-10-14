package dto

import coredto "asset-control/internal/core/dto"

// PrinterRequestDTO filters GET /printers.
type PrinterRequestDTO struct {
	coredto.PaginatedRequestDTO
	BranchID      *uint64 `form:"branch_id"`
	DepartmentsID *uint64 `form:"departments_id"`
}

var PrinterSortableFields = []string{"id", "model", "created_at"}

type PrinterInputDTO struct {
	AssetNumber      uint64  `json:"asset_number" binding:"required,gt=0,lte=9223372036854775807"`
	Model            string  `json:"model" binding:"required,max=100"`
	IP               string  `json:"ip" binding:"required,ip"`
	MachineCondition *string `json:"machine_condition" binding:"omitempty,max=100"`
	BranchID         uint64  `json:"branch_id" binding:"required,gt=0,lte=9223372036854775807"`
	DepartmentsID    uint64  `json:"departments_id" binding:"required,gt=0,lte=9223372036854775807"`
}

func (input PrinterInputDTO) Validated() (PrinterInputDTO, error) {
	err := input.Validate()
	return input, err
}
