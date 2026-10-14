package dto

import (
	"fmt"
	"slices"
	"strings"

	"asset-control/internal/core/exception"
)

type Direction string

const (
	DirectionAsc  Direction = "ASC"
	DirectionDesc Direction = "DESC"
)

// UnmarshalParam uppercases the query value, so ?direction=desc is accepted.
func (d *Direction) UnmarshalParam(value string) error {
	*d = Direction(strings.ToUpper(value))
	return nil
}

// PaginatedRequestDTO is embedded by every XxxRequestDTO. Page is zero-based;
// PageSize 0 lets the resource choose; an empty Direction means ASC.
type PaginatedRequestDTO struct {
	Page      int64     `form:"page" binding:"gte=0"`
	PageSize  int64     `form:"page_size" binding:"gte=0,lte=100"`
	SortBy    string    `form:"sort_by"`
	Direction Direction `form:"direction" binding:"omitempty,oneof=ASC DESC"`
}

// Paginated is satisfied by every request DTO that embeds PaginatedRequestDTO.
type Paginated interface {
	Pagination() PaginatedRequestDTO
}

func (r PaginatedRequestDTO) Pagination() PaginatedRequestDTO { return r }

func (r PaginatedRequestDTO) EffectivePageSize(defaultPageSize int64) int64 {
	if r.PageSize == 0 {
		return defaultPageSize
	}
	return r.PageSize
}

// ValidateSortBy accepts an empty value (default order) or one of the entity's sortable fields.
func (r PaginatedRequestDTO) ValidateSortBy(sortableFields []string) error {
	if r.SortBy == "" || slices.Contains(sortableFields, r.SortBy) {
		return nil
	}
	return exception.NewValidationError(fmt.Sprintf(
		"sort_by must be one of the following values: %s", strings.Join(sortableFields, ", ")))
}
