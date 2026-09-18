package entities

import "time"

type Phone struct {
	ID           uint64    `bson:"id" json:"id"`
	Equipament   *string   `bson:"equipament" json:"equipament"`
	Model        string    `bson:"model" json:"model"`
	IMEI         string    `bson:"imei" json:"imei"`
	SerialNumber string    `bson:"serial_number" json:"serial_number"`
	HasChip      *bool     `bson:"has_chip" json:"has_chip"`
	Situation    string    `bson:"situation" json:"situation"`
	EmployeeID   *uint64   `bson:"employee_id" json:"employee_id"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

func (Phone) CollectionName() string {
	return "phones"
}
