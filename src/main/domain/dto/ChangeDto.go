package dto

import "time"

type ChangeObjectDto struct {
	Before any `json:"before"`
	After  any `json:"after"`
}

type ChangeDto struct {
	Time       time.Time       `json:"time"`
	UserID     string          `json:"user_id"`
	Resource   string          `json:"resource"`
	ResourceID int64           `json:"resource_id"`
	Action     string          `json:"action"`
	Object     ChangeObjectDto `json:"object"`
}
