package dto

import (
	"errors"
	"testing"

	"asset-control/internal/core/exception"
)

func TestNewPaginatedResponseDTO(t *testing.T) {
	for _, tc := range []struct {
		total, page, contentSize int64
		want                     PagingDTO
	}{
		{0, 0, 0, PagingDTO{TotalElements: 0, Size: 20, Page: 0, TotalPages: 0}},
		{42, 0, 20, PagingDTO{TotalElements: 42, Size: 20, Page: 0, TotalPages: 3, HasContent: true, HasNext: true}},
		{42, 2, 2, PagingDTO{TotalElements: 42, Size: 20, Page: 2, TotalPages: 3, HasContent: true, HasPrevious: true}},
		{42, 5, 0, PagingDTO{TotalElements: 42, Size: 20, Page: 5, TotalPages: 3, HasPrevious: true}},
	} {
		var content []int
		for range tc.contentSize {
			content = append(content, 1)
		}
		page := NewPaginatedResponseDTO(content, tc.page, 20, tc.total)
		if page.Paging != tc.want || page.Content == nil {
			t.Errorf("%+v: got %+v", tc, page.Paging)
		}
	}
}

func TestPaginatedRequestDTO(t *testing.T) {
	request := PaginatedRequestDTO{SortBy: "name"}
	if request.EffectivePageSize(20) != 20 || (PaginatedRequestDTO{PageSize: 5}).EffectivePageSize(20) != 5 {
		t.Fatal("page_size 0 must use the resource default")
	}
	if err := request.ValidateSortBy([]string{"id", "name"}); err != nil {
		t.Fatal(err)
	}
	err := request.ValidateSortBy([]string{"id"})
	var validationError *exception.ValidationError
	if !errors.As(err, &validationError) || validationError.Messages[0] != "sort_by must be one of the following values: id" {
		t.Fatalf("got %v", err)
	}
	var direction Direction
	if _ = direction.UnmarshalParam("desc"); direction != DirectionDesc {
		t.Fatal("direction must be case-insensitive")
	}
}

func TestDescribe(t *testing.T) {
	type assignment struct {
		EmployeeIDs []uint64 `form:"employee_id"`
	}
	type machineRequestDTO struct {
		PaginatedRequestDTO
		assignment
		Name   *string `form:"name"`
		Secret string
	}
	name := "Ana"
	got := Describe(machineRequestDTO{PaginatedRequestDTO: PaginatedRequestDTO{PageSize: 5}, assignment: assignment{EmployeeIDs: []uint64{4}}, Name: &name, Secret: "x"})
	if got != "machineRequestDTO(page=0, page_size=5, sort_by=, direction=, employee_id=[4], name=Ana)" {
		t.Fatal(got)
	}
}
