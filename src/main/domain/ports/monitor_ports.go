package ports

import (
	"context"

	"asset-control/src/main/domain/dto"
)

type MonitorPorts interface {
	CreateMonitor(ctx context.Context, input dto.CreateMonitorDto) (dto.MonitorDto, error)
	UpdateMonitor(ctx context.Context, id uint64, input dto.CreateMonitorDto) (dto.MonitorDto, error)
	GetAllMonitors(ctx context.Context) ([]dto.MonitorDto, error)
	FindByID(ctx context.Context, id uint64) (dto.MonitorDto, error)
	FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.MonitorDto, error)
	GetUnassignedMonitors(ctx context.Context) ([]dto.MonitorDto, error)
}
