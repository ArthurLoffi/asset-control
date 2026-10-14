package adapter

import "testing"

func TestEffectiveEmployeeIDs(t *testing.T) {
	id := uint64(4)
	if got := effectiveEmployeeIDs(&id, nil); len(got) != 1 || got[0] != id {
		t.Fatalf("legacy assignment: %v", got)
	}
	if got := effectiveEmployeeIDs(&id, []uint64{}); len(got) != 0 {
		t.Fatalf("empty plural assignment must override legacy value: %v", got)
	}
}

func TestWithoutEmployeeID(t *testing.T) {
	got := withoutEmployeeID([]uint64{2, 3, 4}, 3)
	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Fatalf("removed other assignments: %v", got)
	}
}
