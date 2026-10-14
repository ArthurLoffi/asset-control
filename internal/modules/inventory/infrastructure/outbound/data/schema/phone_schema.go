package schema

import "time"

const PhoneCollection = "phones"

// PhoneSortColumns maps the public sort_by value to the stored field.
var PhoneSortColumns = map[string]string{"id": "id", "model": "model", "created_at": "created_at"}

type PhoneDocument struct {
	ID           uint64    `bson:"id"`
	Equipament   *string   `bson:"equipament"`
	Model        string    `bson:"model"`
	IMEI         string    `bson:"imei"`
	SerialNumber string    `bson:"serial_number"`
	HasChip      *bool     `bson:"has_chip"`
	Situation    string    `bson:"situation"`
	EmployeeID   *uint64   `bson:"employee_id"`
	EmployeeIDs  []uint64  `bson:"employee_ids"`
	CreatedAt    time.Time `bson:"created_at"`
}
