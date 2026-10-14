package adapter

import (
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MonitorRepository struct {
	documents documentRepository[schema.Monitor]
}

var _ port.MonitorPorts = (*MonitorRepository)(nil)

func NewMonitorRepository(db *mongo.Database) *MonitorRepository {
	return &MonitorRepository{documents: newDocumentRepository[schema.Monitor](db, (schema.Monitor{}).CollectionName())}
}

func monitorDocument(input dto.CreateMonitorDto) schema.Monitor {
	return schema.Monitor{
		HasMonitor:       input.HasMonitor,
		AssetNumber:      input.AssetNumber,
		Brand:            input.Brand,
		Model:            input.Model,
		SerialNumber:     input.SerialNumber,
		MonitorCondition: input.MonitorCondition,
		Situation:        input.Situation,
		EmployeeID:       input.EmployeeID,
		EmployeeIDs:      input.EmployeeIDs,
	}
}

func monitorDTO(document schema.Monitor) dto.MonitorDto {
	result := dto.MonitorDto(document)
	result.EmployeeIDs = effectiveEmployeeIDs(result.EmployeeID, result.EmployeeIDs)
	return result
}

func (r *MonitorRepository) CreateMonitor(ctx context.Context, input dto.CreateMonitorDto) (dto.MonitorDto, error) {
	document, err := r.documents.create(ctx, monitorDocument(input))
	return monitorDTO(document), err
}

func (r *MonitorRepository) UpdateMonitor(ctx context.Context, id uint64, input dto.CreateMonitorDto) (dto.MonitorDto, error) {
	document, err := r.documents.update(ctx, id, monitorDocument(input))
	return monitorDTO(document), err
}

func (r *MonitorRepository) GetAllMonitors(ctx context.Context) ([]dto.MonitorDto, error) {
	documents, err := r.documents.find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.MonitorDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, monitorDTO(document))
	}
	return result, nil
}

func (r *MonitorRepository) FindByID(ctx context.Context, id uint64) (dto.MonitorDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"id": id})
	return monitorDTO(document), err
}

func (r *MonitorRepository) FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.MonitorDto, error) {
	documents, err := r.documents.find(ctx, assignedTo(employeeID))
	if err != nil {
		return nil, err
	}
	result := make([]dto.MonitorDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, monitorDTO(document))
	}
	return result, nil
}

func (r *MonitorRepository) GetUnassignedMonitors(ctx context.Context) ([]dto.MonitorDto, error) {
	documents, err := r.documents.find(ctx, unassignedFilter())
	if err != nil {
		return nil, err
	}
	result := make([]dto.MonitorDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, monitorDTO(document))
	}
	return result, nil
}
