package usecase

import (
	"asset-control/internal/core/audit"
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
)

type MachineUseCase struct {
	repository port.MachinePorts
	changes    audit.Recorder
}

func NewMachineUseCase(repository port.MachinePorts, changePorts ...audit.Recorder) *MachineUseCase {
	var changes audit.Recorder
	if len(changePorts) > 0 {
		changes = changePorts[0]
	}
	return &MachineUseCase{repository: repository, changes: changes}
}

func (uc *MachineUseCase) CreateMachine(ctx context.Context, input dto.CreateMachineDto) (dto.MachineDto, error) {
	if err := input.Validate(); err != nil {
		return dto.MachineDto{}, err
	}
	return uc.repository.CreateMachine(ctx, input)
}

func (uc *MachineUseCase) UpdateMachine(ctx context.Context, id uint64, input dto.CreateMachineDto) (dto.MachineDto, error) {
	if err := input.Validate(); err != nil {
		return dto.MachineDto{}, err
	}
	return audit.UpdateWithChange(
		ctx,
		uc.changes,
		"machine",
		int64(id),
		func(txCtx context.Context) (dto.MachineDto, error) { return uc.repository.FindByID(txCtx, id) },
		func(txCtx context.Context) (dto.MachineDto, error) {
			return uc.repository.UpdateMachine(txCtx, id, input)
		},
	)
}

func (uc *MachineUseCase) GetAllMachines(ctx context.Context) ([]dto.MachineDto, error) {
	return uc.repository.GetAllMachines(ctx)
}

func (uc *MachineUseCase) FindByID(ctx context.Context, id uint64) (dto.MachineDto, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *MachineUseCase) FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.MachineDto, error) {
	return uc.repository.FindByEmployeeID(ctx, employeeID)
}

func (uc *MachineUseCase) GetUnassignedMachines(ctx context.Context) ([]dto.MachineDto, error) {
	return uc.repository.GetUnassignedMachines(ctx)
}
