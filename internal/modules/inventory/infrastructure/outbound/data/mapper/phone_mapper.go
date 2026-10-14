package mapper

import (
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
)

func ToPhoneDocument(input dto.PhoneInputDTO) schema.PhoneDocument {
	return schema.PhoneDocument{
		Equipament:   input.Equipament,
		Model:        input.Model,
		IMEI:         input.IMEI,
		SerialNumber: input.SerialNumber,
		HasChip:      input.HasChip,
		Situation:    input.Situation,
		EmployeeID:   input.EmployeeID,
		EmployeeIDs:  input.EmployeeIDs,
	}
}

func ToPhoneResponseDTO(document schema.PhoneDocument) dto.PhoneResponseDTO {
	response := dto.PhoneResponseDTO(document)
	response.EmployeeIDs = effectiveEmployeeIDs(response.EmployeeID, response.EmployeeIDs)
	return response
}
