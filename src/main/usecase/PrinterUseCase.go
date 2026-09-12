package usecase

import (
	"context"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type PrinterUseCase struct{ repository ports.PrinterPorts }

func NewPrinterUseCase(repository ports.PrinterPorts) *PrinterUseCase {
	return &PrinterUseCase{repository: repository}
}

func (uc *PrinterUseCase) CreatePrinter(ctx context.Context, input dto.CreatePrinterDto) (dto.PrinterDto, error) {
	if err := input.Validate(); err != nil {
		return dto.PrinterDto{}, err
	}
	return uc.repository.CreatePrinter(ctx, input)
}

func (uc *PrinterUseCase) UpdatePrinter(ctx context.Context, id uint64, input dto.CreatePrinterDto) (dto.PrinterDto, error) {
	if err := input.Validate(); err != nil {
		return dto.PrinterDto{}, err
	}
	return uc.repository.UpdatePrinter(ctx, id, input)
}

func (uc *PrinterUseCase) GetAllPrinters(ctx context.Context) ([]dto.PrinterDto, error) {
	return uc.repository.GetAllPrinters(ctx)
}

func (uc *PrinterUseCase) FindByID(ctx context.Context, id uint64) (dto.PrinterDto, error) {
	return uc.repository.FindByID(ctx, id)
}
