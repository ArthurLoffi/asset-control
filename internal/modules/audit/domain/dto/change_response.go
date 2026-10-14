package dto

import "time"

type ChangedFieldDTO struct {
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

// ChangeResponseDTO omits audit snapshots; only changed scalar fields leave.
type ChangeResponseDTO struct {
	ID            string            `json:"id"`
	Time          time.Time         `json:"time"`
	UserID        string            `json:"user_id"`
	UserName      string            `json:"user_name"`
	Resource      string            `json:"resource"`
	ResourceID    int64             `json:"resource_id"`
	Action        string            `json:"action"`
	ChangedFields []ChangedFieldDTO `json:"changed_fields"`
}
