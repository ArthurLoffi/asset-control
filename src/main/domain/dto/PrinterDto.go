package dto

import "time"

type PrinterDto struct {
	ID               uint64    `json:"id"`
	AssetNumber      uint64    `json:"asset_number"`
	Model            string    `json:"model"`
	IP               string    `json:"ip"`
	MachineCondition *string   `json:"machine_condition"`
	BranchID         uint64    `json:"branch_id"`
	DepartmentsID    uint64    `json:"departments_id"`
	CreatedAt        time.Time `json:"created_at"`
}
