package dto

import "time"

type EmailResponseDTO struct {
	ID    int32  `json:"id" format:"int32" example:"42"`
	User  string `json:"user" example:"Ana"`
	Email string `json:"email" example:"ana@example.com"`
	// PasswordHash contains the plaintext password as stored in usableEmails.
	PasswordHash string    `json:"password_hash" example:"example-password"`
	CreatedAt    time.Time `json:"created_at" format:"date-time" example:"2026-01-01T00:00:00Z"`
}
