package port

import (
	"asset-control/internal/core/resource"
	"asset-control/internal/modules/organization/domain/dto"
)

// Branches and departments are read-only reference lists.
type (
	BranchDataPort     = resource.ReadPort[dto.BranchRequestDTO, dto.BranchResponseDTO]
	DepartmentDataPort = resource.ReadPort[dto.DepartmentRequestDTO, dto.DepartmentResponseDTO]
)
