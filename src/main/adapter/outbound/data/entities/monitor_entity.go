package entities

import "time"

type Monitor struct {
	ID               uint64    `bson:"id" json:"id"`
	HasMonitor       *bool     `bson:"has_monitor" json:"has_monitor"`
	AssetNumber      *uint64   `bson:"asset_number" json:"asset_number"`
	Brand            *string   `bson:"brand" json:"brand"`
	Model            string    `bson:"model" json:"model"`
	SerialNumber     *string   `bson:"serial_number" json:"serial_number"`
	MonitorCondition *string   `bson:"monitor_condition" json:"monitor_condition"`
	Situation        string    `bson:"situation" json:"situation"`
	EmployeeID       *uint64   `bson:"employee_id" json:"employee_id"`
	EmployeeIDs      []uint64  `bson:"employee_ids" json:"employee_ids"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
}

func (Monitor) CollectionName() string {
	return "monitors"
}
