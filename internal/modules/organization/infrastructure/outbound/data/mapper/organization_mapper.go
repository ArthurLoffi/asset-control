package mapper

import (
	"asset-control/internal/modules/organization/domain/dto"
	"asset-control/internal/modules/organization/infrastructure/outbound/data/schema"
)

func ToBranchResponseDTO(document schema.BranchDocument) dto.BranchResponseDTO {
	return dto.BranchResponseDTO(document)
}

func ToDepartmentResponseDTO(document schema.DepartmentDocument) dto.DepartmentResponseDTO {
	return dto.DepartmentResponseDTO(document)
}
