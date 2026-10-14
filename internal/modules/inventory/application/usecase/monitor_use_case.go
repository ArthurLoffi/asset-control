package usecase

import (
	"asset-control/internal/core/audit"
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
)

type MonitorUseCase struct {
	repository port.MonitorPorts
	changes    audit.Recorder
}

func NewMonitorUseCase(repository port.MonitorPorts, changePorts ...audit.Recorder) *MonitorUseCase {
	var changes audit.Recorder
	if len(changePorts) > 0 {
		changes = changePorts[0]
	}
	return &MonitorUseCase{repository: repository, changes: changes}
}

func (uc *MonitorUseCase) CreateMonitor(ctx context.Context, input dto.CreateMonitorDto) (dto.MonitorDto, error) {
	if err := input.Validate(); err != nil {
		return dto.MonitorDto{}, err
	}
	return uc.repository.CreateMonitor(ctx, input)
}

func (uc *MonitorUseCase) UpdateMonitor(ctx context.Context, id uint64, input dto.CreateMonitorDto) (dto.MonitorDto, error) {
	if err := input.Validate(); err != nil {
		return dto.MonitorDto{}, err
	}
	return audit.UpdateWithChange(
		ctx,
		uc.changes,
		"monitor",
		int64(id),
		func(txCtx context.Context) (dto.MonitorDto, error) { return uc.repository.FindByID(txCtx, id) },
		func(txCtx context.Context) (dto.MonitorDto, error) {
			return uc.repository.UpdateMonitor(txCtx, id, input)
		},
	)
}

func (uc *MonitorUseCase) GetAllMonitors(ctx context.Context) ([]dto.MonitorDto, error) {
	return uc.repository.GetAllMonitors(ctx)
}

func (uc *MonitorUseCase) FindByID(ctx context.Context, id uint64) (dto.MonitorDto, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *MonitorUseCase) FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.MonitorDto, error) {
	return uc.repository.FindByEmployeeID(ctx, employeeID)
}

func (uc *MonitorUseCase) GetUnassignedMonitors(ctx context.Context) ([]dto.MonitorDto, error) {
	return uc.repository.GetUnassignedMonitors(ctx)
}
