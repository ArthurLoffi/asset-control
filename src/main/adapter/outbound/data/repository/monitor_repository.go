package repository

import (
	"context"

	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MonitorRepository struct {
	documents documentRepository[entities.Monitor]
}

var _ ports.MonitorPorts = (*MonitorRepository)(nil)

func NewMonitorRepository(db *mongo.Database) *MonitorRepository {
	return &MonitorRepository{documents: newDocumentRepository[entities.Monitor](db, (entities.Monitor{}).CollectionName())}
}

func monitorDocument(input dto.CreateMonitorDto) entities.Monitor {
	return entities.Monitor{
		HasMonitor:       input.HasMonitor,
		AssetNumber:      input.AssetNumber,
		Brand:            input.Brand,
		Model:            input.Model,
		SerialNumber:     input.SerialNumber,
		MonitorCondition: input.MonitorCondition,
		Situation:        input.Situation,
		EmployeeID:       input.EmployeeID,
	}
}

func (r *MonitorRepository) CreateMonitor(ctx context.Context, input dto.CreateMonitorDto) (dto.MonitorDto, error) {
	document, err := r.documents.create(ctx, monitorDocument(input))
	return dto.MonitorDto(document), err
}

func (r *MonitorRepository) UpdateMonitor(ctx context.Context, id uint64, input dto.CreateMonitorDto) (dto.MonitorDto, error) {
	document, err := r.documents.update(ctx, id, monitorDocument(input))
	return dto.MonitorDto(document), err
}

func (r *MonitorRepository) GetAllMonitors(ctx context.Context) ([]dto.MonitorDto, error) {
	documents, err := r.documents.find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MonitorDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.MonitorDto(document))
	}
	return result, nil
}

func (r *MonitorRepository) FindByID(ctx context.Context, id uint64) (dto.MonitorDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"id": id})
	return dto.MonitorDto(document), err
}

func (r *MonitorRepository) FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.MonitorDto, error) {
	documents, err := r.documents.find(ctx, bson.M{"employee_id": employeeID})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MonitorDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.MonitorDto(document))
	}
	return result, nil
}

func (r *MonitorRepository) GetUnassignedMonitors(ctx context.Context) ([]dto.MonitorDto, error) {
	documents, err := r.documents.find(ctx, bson.M{"employee_id": nil})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MonitorDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.MonitorDto(document))
	}
	return result, nil
}
