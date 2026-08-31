package usecase

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type MonitorUseCase struct {
	monitorPorts ports.MonitorPorts
}

func NewMonitorUseCase(monitorPorts ports.MonitorPorts) *MonitorUseCase {
	return &MonitorUseCase{
		monitorPorts: monitorPorts,
	}
}

func (uc *MonitorUseCase) CreateMonitor(monitor dto.CreateMonitorDto) error {
	return uc.monitorPorts.CreateMonitor(monitor)
}

func (uc *MonitorUseCase) UpdateMonitor(id uint64, monitor dto.CreateMonitorDto) (dto.MonitorDto, error) {
	return uc.monitorPorts.UpdateMonitor(id, monitor)
}

func (uc *MonitorUseCase) GetAllMonitors() ([]dto.MonitorDto, error) {
	return uc.monitorPorts.GetAllMonitors()
}

func (uc *MonitorUseCase) FindByID(id uint64) (dto.MonitorDto, error) {
	return uc.monitorPorts.FindByID(id)
}

func (uc *MonitorUseCase) FindByEmployeeID(employeeId uint64) ([]dto.MonitorDto, error) {
	return uc.monitorPorts.FindByEmployeeID(employeeId)
}
