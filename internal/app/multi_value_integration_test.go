//go:build integration

package app

import (
	"net/http"
	"testing"
)

func TestMultipleContactsAndAssignmentsRoundTrip(t *testing.T) {
	db := testDatabase(t)
	ensureSchema(t, db)
	router := testRouter(t, db)

	employee := map[string]any{
		"name": "Ana", "emails": []string{"ana@example.com", "ana.work@example.com"},
		"phones": []string{"111", "222"}, "branch_id": 1, "department_id": 1,
	}
	created := object(t, callAPI(t, router, http.MethodPost, "/api/employees", employee, http.StatusCreated))
	if created["email"] != "ana@example.com" || created["phone"] != "111" || len(created["emails"].([]any)) != 2 {
		t.Fatalf("employee contacts: %v", created)
	}
	employee["emails"] = []string{"new@example.com", "other@example.com"}
	updated := object(t, callAPI(t, router, http.MethodPut, "/api/employees/1", employee, http.StatusOK))
	if updated["email"] != "new@example.com" || len(updated["emails"].([]any)) != 2 {
		t.Fatalf("updated contacts: %v", updated)
	}

	for _, resource := range []struct {
		name  string
		input map[string]any
	}{
		{"machines", map[string]any{"user": "Ana", "type": "PC", "model": "Desktop"}},
		{"monitors", map[string]any{"model": "Display"}},
		{"phones", map[string]any{"model": "Phone", "imei": "123", "serial_number": "SN"}},
	} {
		resource.input["employee_ids"] = []uint64{11, 12}
		created := object(t, callAPI(t, router, http.MethodPost, "/api/"+resource.name, resource.input, http.StatusCreated))
		if created["employee_id"] != float64(11) || len(created["employee_ids"].([]any)) != 2 {
			t.Fatalf("%s create: %v", resource.name, created)
		}
		found := object(t, callAPI(t, router, http.MethodGet, "/api/"+resource.name+"/1", nil, http.StatusOK))
		if len(found["employee_ids"].([]any)) != 2 {
			t.Fatalf("%s read: %v", resource.name, found)
		}
		resource.input["employee_ids"] = []uint64{12, 13}
		updated := object(t, callAPI(t, router, http.MethodPut, "/api/"+resource.name+"/1", resource.input, http.StatusOK))
		if updated["employee_id"] != float64(12) || len(updated["employee_ids"].([]any)) != 2 {
			t.Fatalf("%s update: %v", resource.name, updated)
		}
		if assigned, _ := page(t, callAPI(t, router, http.MethodGet, "/api/"+resource.name+"?employee_id=13", nil, http.StatusOK)); len(assigned) != 1 {
			t.Fatalf("%s assignment filter: %v", resource.name, assigned)
		}
	}
}
