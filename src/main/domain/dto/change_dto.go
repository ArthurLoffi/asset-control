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

// RecentChangeDto omits audit snapshots from the dashboard feed.
type ChangedFieldDto struct {
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

type RecentChangeDto struct {
	ID            string            `json:"id"`
	Time          time.Time         `json:"time"`
	UserID        string            `json:"user_id"`
	UserName      string            `json:"user_name"`
	Resource      string            `json:"resource"`
	ResourceID    int64             `json:"resource_id"`
	Action        string            `json:"action"`
	ChangedFields []ChangedFieldDto `json:"changed_fields"`
}

type RecentChangesResponseDto struct {
	Status string            `json:"status"`
	Data   []RecentChangeDto `json:"data"`
}
