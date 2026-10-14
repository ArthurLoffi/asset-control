package schema

import "time"

const MachineCollection = "machines"

// MachineSortColumns maps the public sort_by value to the stored field.
var MachineSortColumns = map[string]string{"id": "id", "user": "user", "model": "model", "created_at": "created_at"}

type MachineDocument struct {
	ID               uint64    `bson:"id"`
	User             string    `bson:"user"`
	AssetNumber      *int      `bson:"asset_number"`
	Type             string    `bson:"type"`
	Model            string    `bson:"model"`
	SerialNumber     string    `bson:"serial_number"`
	HasKeyboard      *bool     `bson:"has_keyboard"`
	HasMouse         *bool     `bson:"has_mouse"`
	MachineCondition *string   `bson:"machine_condition"`
	Situation        string    `bson:"situation"`
	EmployeeID       *uint64   `bson:"employee_id"`
	EmployeeIDs      []uint64  `bson:"employee_ids"`
	CreatedAt        time.Time `bson:"created_at"`
}
