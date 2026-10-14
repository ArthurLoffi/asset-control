package adapter

import (
	"reflect"
	"testing"

	"asset-control/internal/modules/inventory/domain/dto"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestWithoutEmployeeID(t *testing.T) {
	legacy := uint64(3)
	for _, tc := range []struct {
		single *uint64
		many   []uint64
		want   []uint64
	}{
		{nil, []uint64{2, 3, 4}, []uint64{2, 4}},
		{&legacy, nil, []uint64{}},
		{&legacy, []uint64{5}, []uint64{5}},
	} {
		if got := withoutEmployeeID(tc.single, tc.many, 3); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got %v, want %v", got, tc.want)
		}
	}
}

func TestAssignmentFilter(t *testing.T) {
	if got := assignmentFilter(dto.AssignmentFilterDTO{}); len(got) != 0 {
		t.Fatalf("no filters must match everything: %v", got)
	}
	got := assignmentFilter(dto.AssignmentFilterDTO{EmployeeIDs: []uint64{1, 2}, IsUnassigned: true})
	conditions, isAnd := got["$and"].(bson.A)
	if !isAnd || len(conditions) != 2 {
		t.Fatalf("filters must combine with $and: %v", got)
	}
}
