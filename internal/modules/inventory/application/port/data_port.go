package port

import (
	"asset-control/internal/core/resource"
	"asset-control/internal/modules/inventory/domain/dto"
)

// Each DataPort is implemented by infrastructure/outbound/data/adapter, so use
// cases depend on these abstractions instead of on MongoDB.
type (
	EmployeeDataPort = resource.DataPort[dto.EmployeeRequestDTO, dto.EmployeeInputDTO, dto.EmployeeResponseDTO]
	MachineDataPort  = resource.DataPort[dto.MachineRequestDTO, dto.MachineInputDTO, dto.MachineResponseDTO]
	MonitorDataPort  = resource.DataPort[dto.MonitorRequestDTO, dto.MonitorInputDTO, dto.MonitorResponseDTO]
	PhoneDataPort    = resource.DataPort[dto.PhoneRequestDTO, dto.PhoneInputDTO, dto.PhoneResponseDTO]
	PrinterDataPort  = resource.DataPort[dto.PrinterRequestDTO, dto.PrinterInputDTO, dto.PrinterResponseDTO]
	EmailDataPort    = resource.DataPort[dto.EmailRequestDTO, dto.EmailInputDTO, dto.EmailResponseDTO]
)
