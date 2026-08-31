package usecase

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type MachineUseCase struct {
	machinePorts ports.MachinePorts
}

func NewMachineUseCase(machinePorts ports.MachinePorts) *MachineUseCase {
	return &MachineUseCase{
		machinePorts: machinePorts,
	}
}

func (uc *MachineUseCase) CreateMachine(machine dto.MachineDto) error {
	return uc.machinePorts.CreateMachine(machine)
}

func (uc *MachineUseCase) UpdateMachine(id uint64, machine dto.MachineDto) (dto.MachineDto, error) {
	return uc.machinePorts.UpdateMachine(id, machine)
}

func (uc *MachineUseCase) GetAllMachines() ([]dto.MachineDto, error) {
	return uc.machinePorts.GetAllMachines()
}

func (uc *MachineUseCase) FindByID(id uint64) (dto.MachineDto, error) {

	return uc.machinePorts.FindByID(id)
}

func (uc *MachineUseCase) FindByEmployeeID(employeeId uint64) ([]dto.MachineDto, error) {
	return uc.machinePorts.FindByEmployeeID(employeeId)
}
