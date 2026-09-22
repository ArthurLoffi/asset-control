package repository

import (
	"asset-control/src/main/adapter/outbound/data/database"
	"asset-control/src/main/domain"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DeactivationRepository struct{ db *mongo.Database }

var _ ports.DeactivationPort = (*DeactivationRepository)(nil)

func NewDeactivationRepository(db *mongo.Database) *DeactivationRepository {
	return &DeactivationRepository{db: db}
}

func (r *DeactivationRepository) collection(resource string) (*mongo.Collection, error) {
	switch resource {
	case "employees", "machines", "phones", "monitors", "printers":
		return r.db.Collection(resource), nil
	case "emails":
		return r.db.Collection(database.UsableEmailsCollection), nil
	default:
		return nil, domain.ErrValidation
	}
}

func (r *DeactivationRepository) Deactivate(ctx context.Context, resource string, id uint64) (dto.DeactivationDto, error) {
	if resource != "employees" {
		return r.deactivateOne(ctx, resource, id)
	}
	session, err := r.db.Client().StartSession()
	if err != nil {
		return dto.DeactivationDto{}, persistenceError(err)
	}
	defer session.EndSession(context.Background())
	value, err := session.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		result, err := r.deactivateOne(txCtx, resource, id)
		if err != nil {
			return nil, err
		}
		for _, collection := range []string{"machines", "monitors", "phones"} {
			_, err := r.db.Collection(collection).UpdateMany(txCtx,
				activeFilter(bson.M{"employee_id": id}),
				bson.M{"$set": bson.M{"employee_id": nil}})
			if err != nil {
				return nil, persistenceError(err)
			}
		}
		return result, nil
	})
	if err != nil {
		return dto.DeactivationDto{}, persistenceError(err)
	}
	return value.(dto.DeactivationDto), nil
}

func (r *DeactivationRepository) deactivateOne(ctx context.Context, resource string, id uint64) (dto.DeactivationDto, error) {
	var result dto.DeactivationDto
	collection, err := r.collection(resource)
	if err != nil {
		return result, err
	}
	// A conditional atomic update records the database time exactly once, even
	// when requests race or a client retries after losing the response.
	err = collection.FindOneAndUpdate(ctx, activeFilter(bson.M{"id": id}),
		bson.M{"$currentDate": bson.M{"deactivated_at": true}},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = collection.FindOne(ctx, bson.M{"id": id, "deactivated_at": bson.M{"$type": "date"}}).Decode(&result)
	}
	return result, persistenceError(err)
}

func (r *DeactivationRepository) ListInactive(ctx context.Context, resource string, after uint64, limit int) (dto.InactivePageDto, error) {
	page := dto.InactivePageDto{Records: make([]dto.DeactivationDto, 0)}
	collection, err := r.collection(resource)
	if err != nil {
		return page, err
	}
	cursor, err := collection.Find(ctx,
		bson.M{"deactivated_at": bson.M{"$type": "date"}, "id": bson.M{"$gt": after}},
		options.Find().SetSort(bson.D{{Key: "id", Value: 1}}).SetLimit(int64(limit+1)).
			SetProjection(bson.M{"_id": 0, "id": 1, "deactivated_at": 1}))
	if err != nil {
		return page, persistenceError(err)
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &page.Records); err != nil {
		return page, persistenceError(err)
	}
	if len(page.Records) > limit {
		page.Records = page.Records[:limit]
		last := page.Records[limit-1].ID
		page.NextAfter = &last
	}
	return page, nil
}
