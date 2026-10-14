package dto

import "time"

type DeactivationResponseDTO struct {
	ID            uint64    `json:"id" bson:"id"`
	DeactivatedAt time.Time `json:"deactivated_at" bson:"deactivated_at"`
}
