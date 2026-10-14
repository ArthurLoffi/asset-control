package usecase

import (
	"context"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/domain/dto"
)

type PrinterUseCase struct {
	repository ports.PrinterPorts
	changes    ports.ChangePorts
}

func NewPrinterUseCase(repository ports.PrinterPorts, changePorts ...ports.ChangePorts) *PrinterUseCase {
	var changes ports.ChangePorts
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
	return updateWithChange(
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
