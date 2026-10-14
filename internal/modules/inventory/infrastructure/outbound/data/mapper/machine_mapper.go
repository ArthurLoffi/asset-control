package mapper

import (
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
)

func ToMachineDocument(input dto.MachineInputDTO) schema.MachineDocument {
	return schema.MachineDocument{
		User:             input.User,
		AssetNumber:      input.AssetNumber,
		Type:             input.Type,
		Model:            input.Model,
		SerialNumber:     input.SerialNumber,
		HasKeyboard:      input.HasKeyboard,
		HasMouse:         input.HasMouse,
		MachineCondition: input.MachineCondition,
		Situation:        input.Situation,
		EmployeeID:       input.EmployeeID,
		EmployeeIDs:      input.EmployeeIDs,
	}
}

func ToMachineResponseDTO(document schema.MachineDocument) dto.MachineResponseDTO {
	response := dto.MachineResponseDTO(document)
	response.EmployeeIDs = effectiveEmployeeIDs(response.EmployeeID, response.EmployeeIDs)
	return response
}
