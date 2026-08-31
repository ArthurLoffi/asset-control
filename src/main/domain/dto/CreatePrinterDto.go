package dto

type CreatePrinterDto struct {
	AssetNumber      uint64  `json:"asset_number" binding:"required"`
	Model            string  `json:"model" binding:"required"`
	IP               string  `json:"ip" binding:"required"`
	MachineCondition *string `json:"machine_condition"`
	BranchID         uint64  `json:"branch_id" binding:"required"`
	DepartmentsID    uint64  `json:"departments_id" binding:"required"`
}
