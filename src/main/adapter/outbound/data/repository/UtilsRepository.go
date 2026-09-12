package repository

import (
	"context"

	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/ports"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UtilsRepository struct{ db *mongo.Database }

var _ ports.UtilsPort = (*UtilsRepository)(nil)

func NewUtilsRepository(db *mongo.Database) *UtilsRepository { return &UtilsRepository{db: db} }

func (r *UtilsRepository) BranchIdToDescription(ctx context.Context, id uint64) (string, error) {
	var branch entities.Branch
	err := r.db.Collection(branch.CollectionName()).FindOne(ctx, bson.M{"id": id}).Decode(&branch)
	return branch.Name, persistenceError(err)
}

func (r *UtilsRepository) DepartmentIdToDescription(ctx context.Context, id uint64) (string, error) {
	var department entities.Department
	err := r.db.Collection(department.CollectionName()).FindOne(ctx, bson.M{"id": id}).Decode(&department)
	return department.Name, persistenceError(err)
}
