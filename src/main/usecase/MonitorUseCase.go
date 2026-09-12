package usecase

import (
	"context"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type MonitorUseCase struct{ repository ports.MonitorPorts }

func NewMonitorUseCase(repository ports.MonitorPorts) *MonitorUseCase {
	return &MonitorUseCase{repository: repository}
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
	return uc.repository.UpdateMonitor(ctx, id, input)
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
