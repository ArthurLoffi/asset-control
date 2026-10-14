package dto

import "testing"

func TestEmployeeMultipleContacts(t *testing.T) {
	input := EmployeeInputDTO{
		Name: "Ana", Emails: []string{" first@example.com ", "second@example.com"},
		Phones: []string{" 123 ", "456"}, BranchID: 1, DepartmentID: 2,
	}
	if err := input.Validate(); err != nil {
		t.Fatal(err)
	}
	if input.Email != "first@example.com" || input.Phone != "123" || len(input.Emails) != 2 || len(input.Phones) != 2 {
		t.Fatalf("contacts not normalized: %+v", input)
	}
	for _, emails := range [][]string{{}, {"invalid"}, {"a@example.com", "A@example.com"}} {
		invalid := EmployeeInputDTO{Name: "Ana", Emails: emails, BranchID: 1, DepartmentID: 2}
		if err := invalid.Validate(); err == nil {
			t.Fatalf("accepted invalid emails: %v", emails)
		}
	}
}

func TestMultipleEmployeeAssignments(t *testing.T) {
	input := MachineInputDTO{User: "Ana", Type: "PC", Model: "PC", EmployeeIDs: []uint64{2, 3}}
	if err := input.Validate(); err != nil {
		t.Fatal(err)
	}
	if input.EmployeeID == nil || *input.EmployeeID != 2 {
		t.Fatalf("primary employee not normalized: %+v", input)
	}
	invalid := MachineInputDTO{User: "Ana", Type: "PC", Model: "PC", EmployeeIDs: []uint64{2, 2}}
	if err := invalid.Validate(); err == nil {
		t.Fatal("accepted duplicate assignments")
	}
}
