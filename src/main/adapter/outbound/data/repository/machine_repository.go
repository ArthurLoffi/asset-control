package repository

import (
	"context"

	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MachineRepository struct {
	documents documentRepository[entities.Machine]
}

var _ ports.MachinePorts = (*MachineRepository)(nil)

func NewMachineRepository(db *mongo.Database) *MachineRepository {
	return &MachineRepository{documents: newDocumentRepository[entities.Machine](db, (entities.Machine{}).CollectionName())}
}

func machineDocument(input dto.CreateMachineDto) entities.Machine {
	return entities.Machine{
		User:             input.User,
		AssetNumber:      input.AssetNumber,
		Type:             input.Type,
		Model:            input.Model,
		SerialNumber:     input.SerialNumber,
		HasKeyboard:      input.HasKeyboard,
		HasMouse:         input.HasMouse,
		MachineCondition: input.MachineCondition,
		Situation:        input.Situation,
		EmployeeID:       input.EmployeeID,
	}
}

func (r *MachineRepository) CreateMachine(ctx context.Context, input dto.CreateMachineDto) (dto.MachineDto, error) {
	document, err := r.documents.create(ctx, machineDocument(input))
	return dto.MachineDto(document), err
}

func (r *MachineRepository) UpdateMachine(ctx context.Context, id uint64, input dto.CreateMachineDto) (dto.MachineDto, error) {
	document, err := r.documents.update(ctx, id, machineDocument(input))
	return dto.MachineDto(document), err
}

func (r *MachineRepository) GetAllMachines(ctx context.Context) ([]dto.MachineDto, error) {
	documents, err := r.documents.find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MachineDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.MachineDto(document))
	}
	return result, nil
}

func (r *MachineRepository) FindByID(ctx context.Context, id uint64) (dto.MachineDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"id": id})
	return dto.MachineDto(document), err
}

func (r *MachineRepository) FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.MachineDto, error) {
	documents, err := r.documents.find(ctx, bson.M{"employee_id": employeeID})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MachineDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.MachineDto(document))
	}
	return result, nil
}
