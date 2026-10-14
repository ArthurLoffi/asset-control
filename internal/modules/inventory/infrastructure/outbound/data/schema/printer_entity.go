package schema

import "time"

type Printer struct {
	ID               uint64    `bson:"id" json:"id"`
	AssetNumber      uint64    `bson:"asset_number" json:"asset_number"`
	Model            string    `bson:"model" json:"model"`
	IP               string    `bson:"ip" json:"ip"`
	MachineCondition *string   `bson:"machine_condition" json:"machine_condition"`
	BranchID         uint64    `bson:"branch_id" json:"branch_id"`
	DepartmentsID    uint64    `bson:"departments_id" json:"departments_id"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
}

func (Printer) CollectionName() string {
	return "printers"
}
