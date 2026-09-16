package dto

import "time"

type DeactivationDto struct {
	ID            uint64    `json:"id" bson:"id"`
	DeactivatedAt time.Time `json:"deactivated_at" bson:"deactivated_at"`
}

type InactivePageDto struct {
	Records   []DeactivationDto `json:"records"`
	NextAfter *uint64           `json:"next_after"`
}

type DeactivationResponseDto struct {
	Status string          `json:"status"`
	Data   DeactivationDto `json:"data"`
}

type InactivePageResponseDto struct {
	Status string          `json:"status"`
	Data   InactivePageDto `json:"data"`
}
