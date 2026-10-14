package port

import (
	"context"

	"asset-control/internal/modules/inventory/domain/dto"
)

type MachinePorts interface {
	CreateMachine(ctx context.Context, input dto.CreateMachineDto) (dto.MachineDto, error)
	UpdateMachine(ctx context.Context, id uint64, input dto.CreateMachineDto) (dto.MachineDto, error)
	GetAllMachines(ctx context.Context) ([]dto.MachineDto, error)
	FindByID(ctx context.Context, id uint64) (dto.MachineDto, error)
	FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.MachineDto, error)
	GetUnassignedMachines(ctx context.Context) ([]dto.MachineDto, error)
}
