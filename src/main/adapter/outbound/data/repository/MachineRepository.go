package repository

import (
	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"errors"

	"gorm.io/gorm"
)

type MachineRepository struct {
	db *gorm.DB
}

func NewMachineRepository(db *gorm.DB) *MachineRepository {
	return &MachineRepository{db: db}
}

func (mr *MachineRepository) CreateMachine(machine dto.MachineDto) error {
	machineEntity := entities.Machine{
		User:             machine.User,
		AssetNumber:      machine.AssetNumber,
		Type:             machine.Type,
		Model:            machine.Model,
		SerialNumber:     machine.SerialNumber,
		HasKeyboard:      machine.HasKeyboard,
		HasMouse:         machine.HasMouse,
		MachineCondition: machine.MachineCondition,
		Situation:        machine.Situation,
		EmployeeID:       machine.EmployeeID,
	}

	if err := mr.db.Create(&machineEntity).Error; err != nil {
		return err
	}

	return nil
}

func (mr *MachineRepository) UpdateMachine(
	id uint64,
	machine dto.MachineDto,
) (dto.MachineDto, error) {
	var machineEntity entities.Machine

	if err := mr.db.First(&machineEntity, id).Error; err != nil {
		return dto.MachineDto{}, err
	}

	machineEntity.User = machine.User
	machineEntity.AssetNumber = machine.AssetNumber
	machineEntity.Type = machine.Type
	machineEntity.Model = machine.Model
	machineEntity.SerialNumber = machine.SerialNumber
	machineEntity.HasKeyboard = machine.HasKeyboard
	machineEntity.HasMouse = machine.HasMouse
	machineEntity.MachineCondition = machine.MachineCondition
	machineEntity.Situation = machine.Situation
	machineEntity.EmployeeID = machine.EmployeeID

	if err := mr.db.Save(&machineEntity).Error; err != nil {
		return dto.MachineDto{}, err
	}

	return dto.MachineDto{
		ID:               machineEntity.ID,
		User:             machineEntity.User,
		AssetNumber:      machineEntity.AssetNumber,
		Type:             machineEntity.Type,
		Model:            machineEntity.Model,
		SerialNumber:     machineEntity.SerialNumber,
		HasKeyboard:      machineEntity.HasKeyboard,
		HasMouse:         machineEntity.HasMouse,
		MachineCondition: machineEntity.MachineCondition,
		Situation:        machineEntity.Situation,
		EmployeeID:       machineEntity.EmployeeID,
		CreatedAt:        machineEntity.CreatedAt,
	}, nil
}

func (mr *MachineRepository) GetAllMachines() ([]dto.MachineDto, error) {
	var machineEntities []entities.Machine

	if err := mr.db.Find(&machineEntities).Error; err != nil {
		return nil, err
	}

	machines := make([]dto.MachineDto, 0, len(machineEntities))

	for _, machineEntity := range machineEntities {
		machines = append(machines, dto.MachineDto{
			ID:               machineEntity.ID,
			User:             machineEntity.User,
			AssetNumber:      machineEntity.AssetNumber,
			Type:             machineEntity.Type,
			Model:            machineEntity.Model,
			SerialNumber:     machineEntity.SerialNumber,
			HasKeyboard:      machineEntity.HasKeyboard,
			HasMouse:         machineEntity.HasMouse,
			MachineCondition: machineEntity.MachineCondition,
			Situation:        machineEntity.Situation,
			EmployeeID:       machineEntity.EmployeeID,
			CreatedAt:        machineEntity.CreatedAt,
		})
	}

	return machines, nil
}

func (mr *MachineRepository) FindByID(id uint64) (dto.MachineDto, error) {
	var machineEntity entities.Machine

	if err := mr.db.First(&machineEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.MachineDto{}, err
		}

		return dto.MachineDto{}, err
	}

	return dto.MachineDto{
		ID:               machineEntity.ID,
		User:             machineEntity.User,
		AssetNumber:      machineEntity.AssetNumber,
		Type:             machineEntity.Type,
		Model:            machineEntity.Model,
		SerialNumber:     machineEntity.SerialNumber,
		HasKeyboard:      machineEntity.HasKeyboard,
		HasMouse:         machineEntity.HasMouse,
		MachineCondition: machineEntity.MachineCondition,
		Situation:        machineEntity.Situation,
		EmployeeID:       machineEntity.EmployeeID,
		CreatedAt:        machineEntity.CreatedAt,
	}, nil
}

func (mr *MachineRepository) FindByEmployeeID(employeeID uint64) ([]dto.MachineDto, error) {
	var machineEntities []entities.Machine

	if err := mr.db.
		Where("employee_id = ?", employeeID).
		Find(&machineEntities).Error; err != nil {
		return nil, err
	}

	machines := make([]dto.MachineDto, 0, len(machineEntities))

	for _, machineEntity := range machineEntities {
		machines = append(machines, dto.MachineDto{
			ID:               machineEntity.ID,
			User:             machineEntity.User,
			AssetNumber:      machineEntity.AssetNumber,
			Type:             machineEntity.Type,
			Model:            machineEntity.Model,
			SerialNumber:     machineEntity.SerialNumber,
			HasKeyboard:      machineEntity.HasKeyboard,
			HasMouse:         machineEntity.HasMouse,
			MachineCondition: machineEntity.MachineCondition,
			Situation:        machineEntity.Situation,
			EmployeeID:       machineEntity.EmployeeID,
			CreatedAt:        machineEntity.CreatedAt,
		})
	}

	return machines, nil
}
