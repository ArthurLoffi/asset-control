package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const CountersCollection = "_counters"
const UsableEmailsCollection = "usableEmails"

func EnsureIndexes(ctx context.Context, db *mongo.Database, name string, indexes ...mongo.IndexModel) error {
	if _, err := db.Collection(name).Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("create indexes for %s: %w", name, err)
	}
	return nil
}

// EnsureLookup prepares a read-only reference collection keyed by a unique id.
func EnsureLookup(ctx context.Context, db *mongo.Database, name string) error {
	return EnsureIndexes(ctx, db, name, mongo.IndexModel{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)})
}

// EnsureResource prepares a Collection: unique id, active and inactive list
// indexes, and the id counter. It preserves imported numeric IDs: $max never
// moves a counter backwards, including when several instances start together.
// Run before accepting writes; duplicate imported IDs fail startup.
func EnsureResource(ctx context.Context, db *mongo.Database, name string, extraIndexes ...mongo.IndexModel) error {
	indexes := append([]mongo.IndexModel{
		{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "deactivated_at", Value: 1}, {Key: "id", Value: 1}}},
		{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetName("inactive_by_id").SetPartialFilterExpression(bson.M{"deactivated_at": bson.M{"$type": "date"}})},
	}, extraIndexes...)
	if err := EnsureIndexes(ctx, db, name, indexes...); err != nil {
		return err
	}
	var last struct {
		ID int64 `bson:"id"`
	}
	err := db.Collection(name).FindOne(ctx, bson.D{}, options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}}).SetProjection(bson.M{"id": 1})).Decode(&last)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("read last ID for %s: %w", name, err)
	}
	counters := db.Collection(CountersCollection)
	_, err = counters.UpdateOne(ctx, bson.M{"_id": name}, bson.M{"$max": bson.M{"seq": last.ID}}, options.UpdateOne().SetUpsert(true))
	// Concurrent first-time upserts can race on the built-in _id index.
	if mongo.IsDuplicateKeyError(err) {
		_, err = counters.UpdateOne(ctx, bson.M{"_id": name}, bson.M{"$max": bson.M{"seq": last.ID}})
	}
	if err != nil {
		return fmt.Errorf("initialize IDs for %s: %w", name, err)
	}
	return nil
}
