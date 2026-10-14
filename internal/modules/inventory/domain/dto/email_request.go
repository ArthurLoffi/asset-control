package dto

import (
	"strings"

	coredto "asset-control/internal/core/dto"
)

// EmailRequestDTO filters GET /emails. User matches partially and case-insensitively.
type EmailRequestDTO struct {
	coredto.PaginatedRequestDTO
	User *string `form:"user"`
}

var EmailSortableFields = []string{"id", "user", "created_at"}

type EmailInputDTO struct {
	User  string `json:"user" binding:"required" example:"Ana"`
	Email string `json:"email" binding:"required,email" example:"ana@example.com"`
	// PasswordHash is stored as plaintext, without trimming or hashing.
	PasswordHash string `json:"password_hash" binding:"required" example:"example-password"`
}

func (input *EmailInputDTO) Validate() error {
	input.User = strings.TrimSpace(input.User)
	input.Email = strings.TrimSpace(input.Email)
	return validateInput(input)
}

func (input EmailInputDTO) Validated() (EmailInputDTO, error) {
	err := input.Validate()
	return input, err
}
