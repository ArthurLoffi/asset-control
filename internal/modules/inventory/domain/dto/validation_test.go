package dto

import (
	"errors"
	"math"
	"testing"

	exception "asset-control/internal/core/exception"
)

func ptr[T any](v T) *T { return &v }

func TestValidation(t *testing.T) {
	tests := []struct {
		name  string
		input interface{ Validate() error }
		valid bool
	}{
		{"valid employee", &CreateEmployeeDto{Name: "Ana", Email: "ana@example.com", BranchID: 1, DepartmentID: 2}, true},
		{"missing employee", &CreateEmployeeDto{}, false},
		{"blank employee name", &CreateEmployeeDto{Name: "  ", Email: "ana@example.com", BranchID: 1, DepartmentID: 2}, false},
		{"invalid email", &CreateEmployeeDto{Name: "Ana", Email: "invalid", BranchID: 1, DepartmentID: 2}, false},
		{"overflow branch", &CreateEmployeeDto{Name: "Ana", Email: "ana@example.com", BranchID: math.MaxUint64, DepartmentID: 2}, false},
		{"valid machine", &CreateMachineDto{User: "Ana", Model: "Laptop", Type: "Notebook", SerialNumber: "SN1", HasMouse: ptr(false)}, true},
		{"invalid machine type", &CreateMachineDto{User: "Ana", Model: "Laptop", Type: "Server", SerialNumber: "SN1"}, false},
		{"invalid condition", &CreateMachineDto{User: "Ana", Model: "Laptop", Type: "PC", SerialNumber: "SN1", MachineCondition: ptr("broken")}, false},
		{"valid monitor", &CreateMonitorDto{Model: "Monitor", HasMonitor: ptr(false)}, true},
		{"zero optional asset number", &CreateMonitorDto{Model: "Monitor", AssetNumber: ptr(uint64(0))}, true},
		{"blank monitor", &CreateMonitorDto{Model: "  "}, false},
		{"invalid situation", &CreateMonitorDto{Model: "Monitor", Situation: "invalid"}, false},
		{"overflow relation", &CreateMonitorDto{Model: "Monitor", EmployeeID: ptr(uint64(math.MaxUint64))}, false},
		{"zero relation", &CreateMonitorDto{Model: "Monitor", EmployeeID: ptr(uint64(0))}, false},
		{"valid phone", &CreatePhoneDto{Model: "Phone", IMEI: "123", SerialNumber: "SN", Equipament: ptr("Só chip")}, true},
		{"invalid equipment", &CreatePhoneDto{Model: "Phone", IMEI: "123", SerialNumber: "SN", Equipament: ptr("unknown")}, false},
		{"valid printer IP", &CreatePrinterDto{Model: "Printer", IP: "192.168.100.100", AssetNumber: 1, BranchID: 1, DepartmentsID: 2}, true},
		{"invalid printer IP", &CreatePrinterDto{Model: "Printer", IP: "999.1.1.1", AssetNumber: 1, BranchID: 1, DepartmentsID: 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.valid && err != nil {
				t.Fatal(err)
			}
			if !tt.valid && !errors.Is(err, exception.ErrValidation) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}

func TestDefaultsAndNullableValues(t *testing.T) {
	input := CreateMachineDto{User: " Ana ", Model: "Laptop", Type: "Notebook", SerialNumber: "SN", HasMouse: ptr(false)}
	if err := input.Validate(); err != nil {
		t.Fatal(err)
	}
	if input.Situation != "Em uso" || input.User != "Ana" {
		t.Fatalf("defaults/normalization not applied: %+v", input)
	}
	if input.HasMouse == nil || *input.HasMouse || input.EmployeeID != nil {
		t.Fatal("nullable values changed")
	}
}
