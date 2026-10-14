package mapper

import (
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
)

func ToEmailDocument(input dto.EmailInputDTO) schema.EmailDocument {
	return schema.EmailDocument{User: input.User, Email: input.Email, PasswordHash: input.PasswordHash}
}

func ToEmailResponseDTO(document schema.EmailDocument) dto.EmailResponseDTO {
	return dto.EmailResponseDTO(document)
}
