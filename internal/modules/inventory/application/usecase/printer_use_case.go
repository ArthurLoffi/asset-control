package usecase

import (
	"asset-control/internal/core/audit"
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
)

type PrinterUseCase struct {
	repository port.PrinterPorts
	changes    audit.Recorder
}

func NewPrinterUseCase(repository port.PrinterPorts, changePorts ...audit.Recorder) *PrinterUseCase {
	var changes audit.Recorder
	if len(changePorts) > 0 {
		changes = changePorts[0]
	}
	return &PrinterUseCase{repository: repository, changes: changes}
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
	return audit.UpdateWithChange(
		ctx,
		uc.changes,
		"printer",
		int64(id),
		func(txCtx context.Context) (dto.PrinterDto, error) { return uc.repository.FindByID(txCtx, id) },
		func(txCtx context.Context) (dto.PrinterDto, error) {
			return uc.repository.UpdatePrinter(txCtx, id, input)
		},
	)
}

func (uc *PrinterUseCase) GetAllPrinters(ctx context.Context) ([]dto.PrinterDto, error) {
	return uc.repository.GetAllPrinters(ctx)
}

func (uc *PrinterUseCase) FindByID(ctx context.Context, id uint64) (dto.PrinterDto, error) {
	return uc.repository.FindByID(ctx, id)
}
