package repository

import (
	"context"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/infrastructure/outbound/data/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UtilsRepository struct{ db *mongo.Database }

var _ ports.UtilsPort = (*UtilsRepository)(nil)

func NewUtilsRepository(db *mongo.Database) *UtilsRepository { return &UtilsRepository{db: db} }

func (r *UtilsRepository) GetAllBranchs(ctx context.Context) ([]dto.BranchDto, error) {
	var branch entities.Branch
	cursor, err := r.db.Collection(branch.CollectionName()).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var branches []entities.Branch
	if err := cursor.All(ctx, &branches); err != nil {
		return nil, persistenceError(err)
	}

	result := make([]dto.BranchDto, 0, len(branches))
	for _, b := range branches {
		result = append(result, dto.BranchDto(b))
	}
	return result, nil
}

func (r *UtilsRepository) BranchIdToDescription(ctx context.Context, id uint64) (string, error) {
	var branch entities.Branch
	err := r.db.Collection(branch.CollectionName()).FindOne(ctx, bson.M{"id": id}).Decode(&branch)
	return branch.Name, persistenceError(err)
}

func (r *UtilsRepository) GetAllDepartments(ctx context.Context) ([]dto.DepartmentDto, error) {
	var department entities.Department
	cursor, err := r.db.Collection(department.CollectionName()).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var departments []entities.Department
	if err := cursor.All(ctx, &departments); err != nil {
		return nil, persistenceError(err)
	}

	result := make([]dto.DepartmentDto, 0, len(departments))
	for _, b := range departments {
		result = append(result, dto.DepartmentDto(b))
	}
	return result, nil
}

func (r *UtilsRepository) DepartmentIdToDescription(ctx context.Context, id uint64) (string, error) {
	var department entities.Department
	err := r.db.Collection(department.CollectionName()).FindOne(ctx, bson.M{"id": id}).Decode(&department)
	return department.Name, persistenceError(err)
}
