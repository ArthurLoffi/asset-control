package port

import (
	"context"

	"asset-control/internal/modules/inventory/domain/dto"
)

type PrinterPorts interface {
	CreatePrinter(ctx context.Context, input dto.CreatePrinterDto) (dto.PrinterDto, error)
	UpdatePrinter(ctx context.Context, id uint64, input dto.CreatePrinterDto) (dto.PrinterDto, error)
	GetAllPrinters(ctx context.Context) ([]dto.PrinterDto, error)
	FindByID(ctx context.Context, id uint64) (dto.PrinterDto, error)
}
