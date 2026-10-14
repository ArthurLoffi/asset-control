package schema

import "time"

const MonitorCollection = "monitors"

// MonitorSortColumns maps the public sort_by value to the stored field.
var MonitorSortColumns = map[string]string{"id": "id", "model": "model", "created_at": "created_at"}

type MonitorDocument struct {
	ID               uint64    `bson:"id"`
	HasMonitor       *bool     `bson:"has_monitor"`
	AssetNumber      *uint64   `bson:"asset_number"`
	Brand            *string   `bson:"brand"`
	Model            string    `bson:"model"`
	SerialNumber     *string   `bson:"serial_number"`
	MonitorCondition *string   `bson:"monitor_condition"`
	Situation        string    `bson:"situation"`
	EmployeeID       *uint64   `bson:"employee_id"`
	EmployeeIDs      []uint64  `bson:"employee_ids"`
	CreatedAt        time.Time `bson:"created_at"`
}
