package dto

import (
	"fmt"
	"strings"

	"asset-control/internal/core/exception"
	"asset-control/internal/core/validation"
)

func validateInput(input any) error { return validation.Struct(input) }

func validateChoice(field, value string, choices ...string) error {
	for _, choice := range choices {
		if value == choice {
			return nil
		}
	}
	return exception.NewValidationError(fmt.Sprintf("invalid %s", field))
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
	if input.Emails == nil {
		input.Emails = []string{input.Email}
	}
	if input.Phones == nil {
		input.Phones = []string{}
		if input.Phone != "" {
			input.Phones = append(input.Phones, input.Phone)
		}
	}
	if len(input.Emails) == 0 {
		return exception.NewValidationError("at least one email is required")
	}
	for index := range input.Emails {
		input.Emails[index] = strings.TrimSpace(input.Emails[index])
	}
	for index := range input.Phones {
		input.Phones[index] = strings.TrimSpace(input.Phones[index])
	}
	if err := uniqueStrings("emails", input.Emails); err != nil {
		return err
	}
	if err := uniqueStrings("phones", input.Phones); err != nil {
		return err
	}
	input.Email = input.Emails[0]
	input.Phone = ""
	if len(input.Phones) > 0 {
		input.Phone = input.Phones[0]
	}
	return validateInput(input)
}

func uniqueStrings(field string, values []string) error {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		key := strings.ToLower(value)
		if value == "" || seen[key] {
			return exception.NewValidationError(fmt.Sprintf("%s must contain distinct nonempty values", field))
		}
		seen[key] = true
	}
	return nil
}

func normalizeEmployeeIDs(single **uint64, many *[]uint64) error {
	if *many == nil {
		*many = []uint64{}
		if *single != nil {
			*many = append(*many, **single)
		}
	}
	seen := make(map[uint64]bool, len(*many))
	for _, id := range *many {
		if seen[id] {
			return exception.NewValidationError("employee_ids must be distinct")
		}
		seen[id] = true
	}
	*single = nil
	if len(*many) > 0 {
		*single = &(*many)[0]
	}
	return nil
}

func (input *CreateMachineDto) Validate() error {
	if err := normalizeEmployeeIDs(&input.EmployeeID, &input.EmployeeIDs); err != nil {
		return err
	}
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
	if err := normalizeEmployeeIDs(&input.EmployeeID, &input.EmployeeIDs); err != nil {
		return err
	}
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
	if err := normalizeEmployeeIDs(&input.EmployeeID, &input.EmployeeIDs); err != nil {
		return err
	}
	input.Model = strings.TrimSpace(input.Model)
	input.IMEI = strings.TrimSpace(input.IMEI)
	input.SerialNumber = strings.TrimSpace(input.SerialNumber)
	requiresIdentification := input.Equipament == nil || (*input.Equipament != "Só chip" && *input.Equipament != "Sem chip")
	if requiresIdentification && (input.Model == "" || input.IMEI == "" || input.SerialNumber == "") {
		return exception.NewValidationError("model, imei and serial_number are required for mobile phone records")
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
