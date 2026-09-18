package repository

import (
	"context"
	"regexp"

	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EmployeeRepository struct {
	documents documentRepository[entities.Employee]
}

var _ ports.EmployeesPorts = (*EmployeeRepository)(nil)

func NewEmployeeRepository(db *mongo.Database) *EmployeeRepository {
	return &EmployeeRepository{documents: newDocumentRepository[entities.Employee](db, (entities.Employee{}).CollectionName())}
}

func employeeDocument(input dto.CreateEmployeeDto) entities.Employee {
	return entities.Employee{
		Name:         input.Name,
		Phone:        input.Phone,
		Email:        input.Email,
		BranchID:     input.BranchID,
		DepartmentID: input.DepartmentID,
	}
}

func (r *EmployeeRepository) CreateEmployee(ctx context.Context, input dto.CreateEmployeeDto) (dto.EmployeeDto, error) {
	document, err := r.documents.create(ctx, employeeDocument(input))
	return dto.EmployeeDto(document), err
}

func (r *EmployeeRepository) UpdateEmployee(ctx context.Context, id uint64, input dto.CreateEmployeeDto) (dto.EmployeeDto, error) {
	document, err := r.documents.update(ctx, id, employeeDocument(input))
	return dto.EmployeeDto(document), err
}

func (r *EmployeeRepository) GetAllEmployees(ctx context.Context) ([]dto.EmployeeDto, error) {
	documents, err := r.documents.find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.EmployeeDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.EmployeeDto(document))
	}
	return result, nil
}

func (r *EmployeeRepository) FindByID(ctx context.Context, id uint64) (dto.EmployeeDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"id": id})
	return dto.EmployeeDto(document), err
}

func (r *EmployeeRepository) FindByName(ctx context.Context, name string) (dto.EmployeeDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"name": bson.Regex{Pattern: regexp.QuoteMeta(name), Options: "i"}})
	return dto.EmployeeDto(document), err
}
