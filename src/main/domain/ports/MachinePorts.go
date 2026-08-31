package ports

import "asset-control/src/main/domain/dto"

type MachinePorts interface {
	CreateMachine(machine dto.MachineDto) error

	UpdateMachine(id uint64, machine dto.MachineDto) (dto.MachineDto, error)

	GetAllMachines() ([]dto.MachineDto, error)

	FindByID(id uint64) (dto.MachineDto, error)

	FindByEmployeeID(employeeId uint64) ([]dto.MachineDto, error)
}
