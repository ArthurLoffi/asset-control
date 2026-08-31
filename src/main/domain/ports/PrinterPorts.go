package ports

import "asset-control/src/main/domain/dto"

type PrinterPorts interface {
	CreatePrinter(printer dto.CreatePrinterDto) error

	UpdatePrinter(id uint64, printer dto.CreatePrinterDto) (dto.PrinterDto, error)

	GetAllPrinters() ([]dto.PrinterDto, error)

	FindByID(id uint64) (dto.PrinterDto, error)
}
