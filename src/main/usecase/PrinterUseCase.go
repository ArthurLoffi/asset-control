package usecase

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type PrinterUseCase struct {
	printerPorts ports.PrinterPorts
}

func NewPrinterUseCase(printerPorts ports.PrinterPorts) *PrinterUseCase {
	return &PrinterUseCase{
		printerPorts: printerPorts,
	}
}

func (uc *PrinterUseCase) CreatePrinter(printer dto.CreatePrinterDto) error {
	return uc.printerPorts.CreatePrinter(printer)
}

func (uc *PrinterUseCase) UpdatePrinter(id uint64, printer dto.CreatePrinterDto) (dto.PrinterDto, error) {
	return uc.printerPorts.UpdatePrinter(id, printer)
}

func (uc *PrinterUseCase) GetAllPrinters() ([]dto.PrinterDto, error) {
	return uc.printerPorts.GetAllPrinters()
}

func (uc *PrinterUseCase) FindByID(id uint64) (dto.PrinterDto, error) {
	return uc.printerPorts.FindByID(id)
}
