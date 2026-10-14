package adapter

import (
	"asset-control/internal/core/database"
	"context"

	"asset-control/internal/modules/organization/application/port"
	"asset-control/internal/modules/organization/domain/dto"
	"asset-control/internal/modules/organization/infrastructure/outbound/data/schema"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UtilsRepository struct{ db *mongo.Database }

var _ port.UtilsPort = (*UtilsRepository)(nil)

func NewUtilsRepository(db *mongo.Database) *UtilsRepository { return &UtilsRepository{db: db} }

func (r *UtilsRepository) GetAllBranchs(ctx context.Context) ([]dto.BranchDto, error) {
	var branch schema.Branch
	cursor, err := r.db.Collection(branch.CollectionName()).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var branches []schema.Branch
	if err := cursor.All(ctx, &branches); err != nil {
		return nil, database.PersistenceError(err)
	}

	result := make([]dto.BranchDto, 0, len(branches))
	for _, b := range branches {
		result = append(result, dto.BranchDto(b))
	}
	return result, nil
}

func (r *UtilsRepository) BranchIdToDescription(ctx context.Context, id uint64) (string, error) {
	var branch schema.Branch
	err := r.db.Collection(branch.CollectionName()).FindOne(ctx, bson.M{"id": id}).Decode(&branch)
	return branch.Name, database.PersistenceError(err)
}

func (r *UtilsRepository) GetAllDepartments(ctx context.Context) ([]dto.DepartmentDto, error) {
	var department schema.Department
	cursor, err := r.db.Collection(department.CollectionName()).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var departments []schema.Department
	if err := cursor.All(ctx, &departments); err != nil {
		return nil, database.PersistenceError(err)
	}

	result := make([]dto.DepartmentDto, 0, len(departments))
	for _, b := range departments {
		result = append(result, dto.DepartmentDto(b))
	}
	return result, nil
}

func (r *UtilsRepository) DepartmentIdToDescription(ctx context.Context, id uint64) (string, error) {
	var department schema.Department
	err := r.db.Collection(department.CollectionName()).FindOne(ctx, bson.M{"id": id}).Decode(&department)
	return department.Name, database.PersistenceError(err)
}
