package dto

import coredto "asset-control/internal/core/dto"

// ChangeRequestDTO pages GET /changes, ordered by time with the record id as
// tie-breaker. Pass direction=DESC for the newest changes first.
type ChangeRequestDTO struct {
	coredto.PaginatedRequestDTO
}

var ChangeSortableFields = []string{"time"}
