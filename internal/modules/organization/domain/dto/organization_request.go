package dto

import coredto "asset-control/internal/core/dto"

// BranchRequestDTO filters GET /branches. Name matches partially and case-insensitively.
type BranchRequestDTO struct {
	coredto.PaginatedRequestDTO
	Name *string `form:"name"`
}

// DepartmentRequestDTO filters GET /departments. Name matches partially and case-insensitively.
type DepartmentRequestDTO struct {
	coredto.PaginatedRequestDTO
	Name *string `form:"name"`
}

var OrganizationSortableFields = []string{"id", "name"}
