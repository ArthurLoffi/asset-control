package dto

import (
	"errors"
	"strings"
	"testing"

	"asset-control/src/main/domain"
)

func TestPhoneIdentificationRequirements(t *testing.T) {
	for _, equipment := range []*string{nil, ptr("Celular"), ptr("Só chip"), ptr("Sem chip")} {
		name := "unspecified"
		if equipment != nil {
			name = *equipment
		}
		t.Run(name, func(t *testing.T) {
			optional := name == "Só chip" || name == "Sem chip"
			for _, field := range []string{"model", "imei", "serial_number"} {
				t.Run(field, func(t *testing.T) {
					input := CreatePhoneDto{Equipament: equipment, Model: "Device", IMEI: "012345678901234", SerialNumber: "SN-1"}
					switch field {
					case "model":
						input.Model = "  "
					case "imei":
						input.IMEI = "  "
					case "serial_number":
						input.SerialNumber = "  "
					}
					err := input.Validate()
					if optional && err != nil {
						t.Fatal(err)
					}
					if !optional && !errors.Is(err, domain.ErrValidation) {
						t.Fatalf("expected validation error, got %v", err)
					}
				})
			}
		})
	}
}

func TestPhoneOptionalIdentificationPreservesLimitsAndValues(t *testing.T) {
	for _, equipment := range []string{"Só chip", "Sem chip"} {
		input := CreatePhoneDto{Equipament: ptr(equipment), Model: "Device", IMEI: "012345678901234", SerialNumber: "SN-1"}
		if err := input.Validate(); err != nil {
			t.Fatal(err)
		}
		if input.Model != "Device" || input.IMEI != "012345678901234" || input.SerialNumber != "SN-1" {
			t.Fatal("identification was changed")
		}
		for _, field := range []string{"model", "imei", "serial_number"} {
			oversized := CreatePhoneDto{Equipament: ptr(equipment)}
			switch field {
			case "model":
				oversized.Model = strings.Repeat("a", 101)
			case "imei":
				oversized.IMEI = strings.Repeat("1", 101)
			case "serial_number":
				oversized.SerialNumber = strings.Repeat("a", 101)
			}
			if err := oversized.Validate(); !errors.Is(err, domain.ErrValidation) {
				t.Fatalf("%s: expected length validation error, got %v", field, err)
			}
		}
	}
}
