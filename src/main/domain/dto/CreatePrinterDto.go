package dto

type CreatePrinterDto struct {
	AssetNumber      uint64  `json:"asset_number" binding:"required,gt=0,lte=9223372036854775807"`
	Model            string  `json:"model" binding:"required,max=100"`
	IP               string  `json:"ip" binding:"required,ip"`
	MachineCondition *string `json:"machine_condition" binding:"omitempty,max=100"`
	BranchID         uint64  `json:"branch_id" binding:"required,gt=0,lte=9223372036854775807"`
	DepartmentsID    uint64  `json:"departments_id" binding:"required,gt=0,lte=9223372036854775807"`
}
