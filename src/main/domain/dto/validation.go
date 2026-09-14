package dto

import (
	"fmt"
	"strings"

	"asset-control/src/main/domain"

	"github.com/go-playground/validator/v10"
)

var inputValidator = func() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	return v
}()

func validateInput(input any) error {
	if err := inputValidator.Struct(input); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrValidation, err)
	}
	return nil
}

func validateChoice(field, value string, choices ...string) error {
	for _, choice := range choices {
		if value == choice {
			return nil
		}
	}
	return fmt.Errorf("%w: invalid %s", domain.ErrValidation, field)
}

func validateCondition(field string, value *string) error {
	if value == nil {
		return nil
	}
	return validateChoice(field, *value, "Novo", "Bom", "Ruim")
}

// SQL supplied this default and enforced these enum values. Preserve both in
// application validation now that writes go directly to MongoDB.
func validateSituation(value *string) error {
	if *value == "" {
		*value = "Em uso"
	}
	return validateChoice("situation", *value, "Em uso", "Sem uso", "Com defeito", "Uso temporário", "Manutenção")
}

func (input *CreateEmployeeDto) Validate() error {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	return validateInput(input)
}

func (input *CreateMachineDto) Validate() error {
	input.User = strings.TrimSpace(input.User)
	input.Model = strings.TrimSpace(input.Model)
	input.SerialNumber = strings.TrimSpace(input.SerialNumber)
	if err := validateInput(input); err != nil {
		return err
	}
	if err := validateChoice("type", input.Type, "Notebook", "PC"); err != nil {
		return err
	}
	if err := validateCondition("machine_condition", input.MachineCondition); err != nil {
		return err
	}
	return validateSituation(&input.Situation)
}

func (input *CreateMonitorDto) Validate() error {
	input.Model = strings.TrimSpace(input.Model)
	if err := validateInput(input); err != nil {
		return err
	}
	if err := validateCondition("monitor_condition", input.MonitorCondition); err != nil {
		return err
	}
	return validateSituation(&input.Situation)
}

func (input *CreatePhoneDto) Validate() error {
	input.Model = strings.TrimSpace(input.Model)
	input.IMEI = strings.TrimSpace(input.IMEI)
	input.SerialNumber = strings.TrimSpace(input.SerialNumber)
	requiresIdentification := input.Equipament == nil || (*input.Equipament != "Só chip" && *input.Equipament != "Sem chip")
	if requiresIdentification && (input.Model == "" || input.IMEI == "" || input.SerialNumber == "") {
		return fmt.Errorf("%w: model, imei and serial_number are required for mobile phone records", domain.ErrValidation)
	}
	if err := validateInput(input); err != nil {
		return err
	}
	if input.Equipament != nil {
		if err := validateChoice("equipament", *input.Equipament, "Celular", "Só chip", "Sem chip"); err != nil {
			return err
		}
	}
	return validateSituation(&input.Situation)
}

func (input *CreatePrinterDto) Validate() error {
	input.Model = strings.TrimSpace(input.Model)
	if err := validateInput(input); err != nil {
		return err
	}
	return validateCondition("machine_condition", input.MachineCondition)
}
