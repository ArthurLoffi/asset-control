package mapper

import (
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
)

func ToPrinterDocument(input dto.PrinterInputDTO) schema.PrinterDocument {
	return schema.PrinterDocument{
		AssetNumber:      input.AssetNumber,
		Model:            input.Model,
		IP:               input.IP,
		MachineCondition: input.MachineCondition,
		BranchID:         input.BranchID,
		DepartmentsID:    input.DepartmentsID,
	}
}

func ToPrinterResponseDTO(document schema.PrinterDocument) dto.PrinterResponseDTO {
	return dto.PrinterResponseDTO(document)
}
