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

// EnsureSchema preserves imported numeric IDs. $max never moves a counter
// backwards, including when multiple application instances start together.
// Run before accepting writes; duplicate imported IDs fail startup.
func EnsureSchema(ctx context.Context, db *mongo.Database) error {
	for _, name := range []string{"employees", "machines", "monitors", "phones", "printers", "branches", "departments", UsableEmailsCollection} {
		collection := db.Collection(name)
		indexes := []mongo.IndexModel{{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)}}
		if name != "branches" && name != "departments" {
			indexes = append(indexes,
				mongo.IndexModel{Keys: bson.D{{Key: "deactivated_at", Value: 1}, {Key: "id", Value: 1}}},
				mongo.IndexModel{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetName("inactive_by_id").SetPartialFilterExpression(bson.M{"deactivated_at": bson.M{"$type": "date"}})},
			)
		}
		if name == "machines" || name == "monitors" || name == "phones" {
			indexes = append(indexes, mongo.IndexModel{Keys: bson.D{{Key: "employee_id", Value: 1}}})
			indexes = append(indexes, mongo.IndexModel{Keys: bson.D{{Key: "employee_id", Value: 1}, {Key: "deactivated_at", Value: 1}, {Key: "id", Value: 1}}})
		}
		if _, err := collection.Indexes().CreateMany(ctx, indexes); err != nil {
			return fmt.Errorf("create indexes for %s: %w", name, err)
		}
		if name == "branches" || name == "departments" {
			continue
		}
		var last struct {
			ID int64 `bson:"id"`
		}
		err := collection.FindOne(ctx, bson.D{}, options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}}).SetProjection(bson.M{"id": 1})).Decode(&last)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("read last ID for %s: %w", name, err)
		}
		_, err = db.Collection(CountersCollection).UpdateOne(ctx, bson.M{"_id": name}, bson.M{"$max": bson.M{"seq": last.ID}}, options.UpdateOne().SetUpsert(true))
		// Concurrent first-time upserts can race on the built-in _id index.
		if mongo.IsDuplicateKeyError(err) {
			_, err = db.Collection(CountersCollection).UpdateOne(ctx, bson.M{"_id": name}, bson.M{"$max": bson.M{"seq": last.ID}})
		}
		if err != nil {
			return fmt.Errorf("initialize IDs for %s: %w", name, err)
		}
	}
	return nil
}
