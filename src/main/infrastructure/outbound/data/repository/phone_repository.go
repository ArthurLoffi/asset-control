package repository

import (
	"context"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/infrastructure/outbound/data/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PhoneRepository struct {
	documents documentRepository[entities.Phone]
}

var _ ports.PhonePorts = (*PhoneRepository)(nil)

func NewPhoneRepository(db *mongo.Database) *PhoneRepository {
	return &PhoneRepository{documents: newDocumentRepository[entities.Phone](db, (entities.Phone{}).CollectionName())}
}

func phoneDocument(input dto.CreatePhoneDto) entities.Phone {
	return entities.Phone{
		Equipament:   input.Equipament,
		Model:        input.Model,
		IMEI:         input.IMEI,
		SerialNumber: input.SerialNumber,
		HasChip:      input.HasChip,
		Situation:    input.Situation,
		EmployeeID:   input.EmployeeID,
		EmployeeIDs:  input.EmployeeIDs,
	}
}

func phoneDTO(document entities.Phone) dto.PhoneDto {
	result := dto.PhoneDto(document)
	result.EmployeeIDs = effectiveEmployeeIDs(result.EmployeeID, result.EmployeeIDs)
	return result
}

func (r *PhoneRepository) CreatePhone(ctx context.Context, input dto.CreatePhoneDto) (dto.PhoneDto, error) {
	document, err := r.documents.create(ctx, phoneDocument(input))
	return phoneDTO(document), err
}

func (r *PhoneRepository) UpdatePhone(ctx context.Context, id uint64, input dto.CreatePhoneDto) (dto.PhoneDto, error) {
	document, err := r.documents.update(ctx, id, phoneDocument(input))
	return phoneDTO(document), err
}

func (r *PhoneRepository) GetAllPhones(ctx context.Context) ([]dto.PhoneDto, error) {
	documents, err := r.documents.find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.PhoneDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, phoneDTO(document))
	}
	return result, nil
}

func (r *PhoneRepository) FindByID(ctx context.Context, id uint64) (dto.PhoneDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"id": id})
	return phoneDTO(document), err
}

func (r *PhoneRepository) FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.PhoneDto, error) {
	documents, err := r.documents.find(ctx, assignedTo(employeeID))
	if err != nil {
		return nil, err
	}
	result := make([]dto.PhoneDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, phoneDTO(document))
	}
	return result, nil
}
