package exception

import "asset-control/internal/core/exception"

const (
	EmployeeNotFoundByIDCode = "EMPLOYEE_001"
	MachineNotFoundByIDCode  = "MACHINE_001"
	MonitorNotFoundByIDCode  = "MONITOR_001"
	PhoneNotFoundByIDCode    = "PHONE_001"
	PrinterNotFoundByIDCode  = "PRINTER_001"
	EmailNotFoundByIDCode    = "EMAIL_001"
)

func NewEmployeeNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Employee", EmployeeNotFoundByIDCode, id)
}

func NewMachineNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Machine", MachineNotFoundByIDCode, id)
}

func NewMonitorNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Monitor", MonitorNotFoundByIDCode, id)
}

func NewPhoneNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Phone", PhoneNotFoundByIDCode, id)
}

func NewPrinterNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Printer", PrinterNotFoundByIDCode, id)
}

func NewEmailNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Email", EmailNotFoundByIDCode, id)
}
