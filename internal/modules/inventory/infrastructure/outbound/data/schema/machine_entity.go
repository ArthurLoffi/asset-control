package schema

import "time"

type Machine struct {
	ID               uint64    `bson:"id" json:"id"`
	User             string    `bson:"user" json:"user"`
	AssetNumber      *int      `bson:"asset_number" json:"asset_number"`
	Type             string    `bson:"type" json:"type"`
	Model            string    `bson:"model" json:"model"`
	SerialNumber     string    `bson:"serial_number" json:"serial_number"`
	HasKeyboard      *bool     `bson:"has_keyboard" json:"has_keyboard"`
	HasMouse         *bool     `bson:"has_mouse" json:"has_mouse"`
	MachineCondition *string   `bson:"machine_condition" json:"machine_condition"`
	Situation        string    `bson:"situation" json:"situation"`
	EmployeeID       *uint64   `bson:"employee_id" json:"employee_id"`
	EmployeeIDs      []uint64  `bson:"employee_ids" json:"employee_ids"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
}

func (Machine) CollectionName() string {
	return "machines"
}
