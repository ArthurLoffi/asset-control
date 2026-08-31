package ports

import "asset-control/src/main/domain/dto"

type MonitorPorts interface {
	CreateMonitor(monitor dto.CreateMonitorDto) error

	UpdateMonitor(id uint64, monitor dto.CreateMonitorDto) (dto.MonitorDto, error)

	GetAllMonitors() ([]dto.MonitorDto, error)

	FindByID(id uint64) (dto.MonitorDto, error)

	FindByEmployeeID(employeeId uint64) ([]dto.MonitorDto, error)
}
