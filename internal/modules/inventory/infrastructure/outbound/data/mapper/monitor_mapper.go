package mapper

import (
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
)

func ToMonitorDocument(input dto.MonitorInputDTO) schema.MonitorDocument {
	return schema.MonitorDocument{
		HasMonitor:       input.HasMonitor,
		AssetNumber:      input.AssetNumber,
		Brand:            input.Brand,
		Model:            input.Model,
		SerialNumber:     input.SerialNumber,
		MonitorCondition: input.MonitorCondition,
		Situation:        input.Situation,
		EmployeeID:       input.EmployeeID,
		EmployeeIDs:      input.EmployeeIDs,
	}
}

func ToMonitorResponseDTO(document schema.MonitorDocument) dto.MonitorResponseDTO {
	response := dto.MonitorResponseDTO(document)
	response.EmployeeIDs = effectiveEmployeeIDs(response.EmployeeID, response.EmployeeIDs)
	return response
}
