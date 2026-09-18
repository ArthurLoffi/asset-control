package usecase

import (
	"context"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type MachineUseCase struct {
	repository ports.MachinePorts
	changes    ports.ChangePorts
}

func NewMachineUseCase(repository ports.MachinePorts, changePorts ...ports.ChangePorts) *MachineUseCase {
	var changes ports.ChangePorts
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
	return updateWithChange(
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
