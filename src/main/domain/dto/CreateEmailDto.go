package dto

import "strings"

type CreateEmailDto struct {
	User  string `json:"user" binding:"required" example:"Ana"`
	Email string `json:"email" binding:"required,email" example:"ana@example.com"`
	// PasswordHash is stored as plaintext, without trimming or hashing.
	PasswordHash string `json:"password_hash" binding:"required" example:"example-password"`
}

func (input *CreateEmailDto) Validate() error {
	input.User = strings.TrimSpace(input.User)
	input.Email = strings.TrimSpace(input.Email)
	return validateInput(input)
}
