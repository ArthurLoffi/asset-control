package adapter

import (
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"
	"testing"
)

func TestEmployeeDTOIncludesLegacyContacts(t *testing.T) {
	result := employeeDTO(schema.Employee{Email: "ana@example.com", Phone: "123"})
	if len(result.Emails) != 1 || result.Emails[0] != "ana@example.com" || len(result.Phones) != 1 || result.Phones[0] != "123" {
		t.Fatalf("legacy contacts missing: %+v", result)
	}

	result = employeeDTO(schema.Employee{})
	if len(result.Emails) != 0 || len(result.Phones) != 0 {
		t.Fatalf("empty legacy contacts: %+v", result)
	}
}
