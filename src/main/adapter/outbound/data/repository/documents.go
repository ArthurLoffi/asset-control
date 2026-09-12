package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"asset-control/src/main/adapter/outbound/data/database"
	"asset-control/src/main/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type documentRepository[T any] struct {
	collection *mongo.Collection
	counters   *mongo.Collection
}

func newDocumentRepository[T any](db *mongo.Database, name string) documentRepository[T] {
	return documentRepository[T]{collection: db.Collection(name), counters: db.Collection(database.CountersCollection)}
}

func persistenceError(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return domain.ErrNotFound
	case mongo.IsDuplicateKeyError(err):
		return domain.ErrConflict
	default:
		return err
	}
}

func (r documentRepository[T]) create(ctx context.Context, document T) (T, error) {
	var zero T
	fields, err := documentFields(document)
	if err != nil {
		return zero, err
	}
	var counter struct {
		Seq int64 `bson:"seq"`
	}
	err = r.counters.FindOneAndUpdate(ctx,
		bson.M{"_id": r.collection.Name(), "seq": bson.M{"$lt": int64(math.MaxInt64)}},
		bson.M{"$inc": bson.M{"seq": int64(1)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&counter)
	if err != nil {
		return zero, fmt.Errorf("allocate %s ID (initialize schema before writes): %w", r.collection.Name(), err)
	}
	fields["id"] = counter.Seq
	fields["created_at"] = time.Now().UTC().Truncate(time.Millisecond)
	// Decode before inserting so conversion errors cannot report a failed write.
	encoded, err := bson.Marshal(fields)
	if err != nil {
		return zero, err
	}
	var result T
	if err := bson.Unmarshal(encoded, &result); err != nil {
		return zero, err
	}
	if _, err := r.collection.InsertOne(ctx, fields); err != nil {
		return zero, persistenceError(err)
	}
	return result, nil
}

// $set retains unknown imported fields and created_at, and writes explicit
// nulls/false values instead of omitting them. It never upserts missing records.
func (r documentRepository[T]) update(ctx context.Context, id uint64, document T) (T, error) {
	var result T
	fields, err := documentFields(document)
	if err != nil {
		return result, err
	}
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"id": id}, bson.M{"$set": fields}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
	return result, persistenceError(err)
}

func documentFields(document any) (bson.M, error) {
	encoded, err := bson.Marshal(document)
	if err != nil {
		return nil, err
	}
	var fields bson.M
	if err := bson.Unmarshal(encoded, &fields); err != nil {
		return nil, err
	}
	delete(fields, "_id")
	delete(fields, "id")
	delete(fields, "created_at")
	return fields, nil
}

func (r documentRepository[T]) findOne(ctx context.Context, filter bson.M) (T, error) {
	var result T
	err := r.collection.FindOne(ctx, filter, options.FindOne().SetSort(bson.D{{Key: "id", Value: 1}})).Decode(&result)
	return result, persistenceError(err)
}

func (r documentRepository[T]) find(ctx context.Context, filter bson.M) ([]T, error) {
	result := make([]T, 0)
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "id", Value: 1}}))
	if err != nil {
		return nil, persistenceError(err)
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &result); err != nil {
		return nil, persistenceError(err)
	}
	return result, nil
}
