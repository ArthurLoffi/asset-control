package adapter

import (
	"context"
	"fmt"
	"math"
	"time"

	"asset-control/internal/core/database"

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

func (r documentRepository[T]) create(ctx context.Context, document T) (T, error) {
	return r.createWithIDLimit(ctx, document, math.MaxInt64)
}

// Email IDs retain their imported BSON int32 representation and cannot overflow.
func (r documentRepository[T]) createWithIDLimit(ctx context.Context, document T, maxID int64) (T, error) {
	var zero T
	fields, err := documentFields(document)
	if err != nil {
		return zero, err
	}
	var counter struct {
		Seq int64 `bson:"seq"`
	}
	err = r.counters.FindOneAndUpdate(ctx,
		bson.M{"_id": r.collection.Name(), "seq": bson.M{"$lt": maxID}},
		bson.M{"$inc": bson.M{"seq": int64(1)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&counter)
	if err != nil {
		return zero, fmt.Errorf("allocate %s ID (initialize schema before writes): %w", r.collection.Name(), err)
	}
	fields["id"] = counter.Seq
	if maxID == math.MaxInt32 {
		fields["id"] = int32(counter.Seq)
	}
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
		return zero, database.PersistenceError(err)
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
	err = r.collection.FindOneAndUpdate(ctx, activeFilter(bson.M{"id": id}), bson.M{"$set": fields}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
	return result, database.PersistenceError(err)
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
	delete(fields, "deactivated_at")
	return fields, nil
}

// Null equality includes imported records without lifecycle metadata. Copy the
// filter so callers cannot accidentally retain or override the active scope.
func activeFilter(filter bson.M) bson.M {
	result := make(bson.M, len(filter)+1)
	for key, value := range filter {
		result[key] = value
	}
	result["deactivated_at"] = nil
	return result
}

func (r documentRepository[T]) findOne(ctx context.Context, filter bson.M) (T, error) {
	var result T
	err := r.collection.FindOne(ctx, activeFilter(filter), options.FindOne().SetSort(bson.D{{Key: "id", Value: 1}})).Decode(&result)
	return result, database.PersistenceError(err)
}

func (r documentRepository[T]) find(ctx context.Context, filter bson.M) ([]T, error) {
	result := make([]T, 0)
	cursor, err := r.collection.Find(ctx, activeFilter(filter), options.Find().SetSort(bson.D{{Key: "id", Value: 1}}))
	if err != nil {
		return nil, database.PersistenceError(err)
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &result); err != nil {
		return nil, database.PersistenceError(err)
	}
	return result, nil
}
