package mapper

import (
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
)

func ToEmployeeDocument(input dto.EmployeeInputDTO) schema.EmployeeDocument {
	return schema.EmployeeDocument{
		Name:         input.Name,
		Phone:        input.Phone,
		Phones:       input.Phones,
		Email:        input.Email,
		Emails:       input.Emails,
		BranchID:     input.BranchID,
		DepartmentID: input.DepartmentID,
	}
}

// ToEmployeeResponseDTO fills the contact lists of imported records that only
// have the singular email and phone.
func ToEmployeeResponseDTO(document schema.EmployeeDocument) dto.EmployeeResponseDTO {
	response := dto.EmployeeResponseDTO(document)
	if response.Emails == nil {
		response.Emails = nonEmpty(response.Email)
	}
	if response.Phones == nil {
		response.Phones = nonEmpty(response.Phone)
	}
	return response
}

func nonEmpty(value string) []string {
	if value == "" {
		return []string{}
	}
	return []string{value}
}
