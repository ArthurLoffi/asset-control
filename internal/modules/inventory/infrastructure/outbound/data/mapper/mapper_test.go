package mapper

import (
	"testing"

	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
)

func TestEffectiveEmployeeIDs(t *testing.T) {
	id := uint64(4)
	if got := effectiveEmployeeIDs(&id, nil); len(got) != 1 || got[0] != id {
		t.Fatalf("legacy assignment: %v", got)
	}
	if got := effectiveEmployeeIDs(&id, []uint64{}); len(got) != 0 {
		t.Fatalf("empty plural assignment must override legacy value: %v", got)
	}
	if got := ToMachineResponseDTO(schema.MachineDocument{EmployeeID: &id}); len(got.EmployeeIDs) != 1 {
		t.Fatalf("machine mapper ignored legacy assignment: %+v", got)
	}
}

func TestEmployeeResponseIncludesLegacyContacts(t *testing.T) {
	response := ToEmployeeResponseDTO(schema.EmployeeDocument{Email: "ana@example.com", Phone: "123"})
	if len(response.Emails) != 1 || response.Emails[0] != "ana@example.com" || len(response.Phones) != 1 || response.Phones[0] != "123" {
		t.Fatalf("legacy contacts missing: %+v", response)
	}
	response = ToEmployeeResponseDTO(schema.EmployeeDocument{})
	if response.Emails == nil || len(response.Emails) != 0 || response.Phones == nil || len(response.Phones) != 0 {
		t.Fatalf("empty legacy contacts: %+v", response)
	}
}
