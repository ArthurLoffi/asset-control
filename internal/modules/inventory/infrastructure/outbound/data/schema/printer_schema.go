package schema

import "time"

const PrinterCollection = "printers"

// PrinterSortColumns maps the public sort_by value to the stored field.
var PrinterSortColumns = map[string]string{"id": "id", "model": "model", "created_at": "created_at"}

type PrinterDocument struct {
	ID               uint64    `bson:"id"`
	AssetNumber      uint64    `bson:"asset_number"`
	Model            string    `bson:"model"`
	IP               string    `bson:"ip"`
	MachineCondition *string   `bson:"machine_condition"`
	BranchID         uint64    `bson:"branch_id"`
	DepartmentsID    uint64    `bson:"departments_id"`
	CreatedAt        time.Time `bson:"created_at"`
}
