package audit

import (
	"context"
	"time"
)

type ChangeObject struct {
	Before any `json:"before"`
	After  any `json:"after"`
}

type Change struct {
	Time       time.Time    `json:"time"`
	UserID     string       `json:"user_id"`
	Resource   string       `json:"resource"`
	ResourceID int64        `json:"resource_id"`
	Action     string       `json:"action"`
	Object     ChangeObject `json:"object"`
}

// Recorder is implemented by the audit module, so other modules can record
// changes without importing it.
type Recorder interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
	CreateChange(ctx context.Context, change Change) error
}
