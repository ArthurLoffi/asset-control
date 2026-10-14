package repository

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/infrastructure/outbound/data/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChangeRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewChangeRepository(db *mongo.Database) *ChangeRepository {
	return &ChangeRepository{
		db:         db,
		collection: db.Collection((entities.Change{}).CollectionName()),
	}
}

func (r *ChangeRepository) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		return nil, fn(txCtx)
	})
	return err
}

func (r *ChangeRepository) CreateChange(ctx context.Context, input dto.ChangeDto) error {
	document := entities.Change{
		Time:       input.Time,
		UserID:     input.UserID,
		Resource:   input.Resource,
		ResourceID: input.ResourceID,
		Action:     input.Action,
		Object: entities.ChangeObject{
			Before: auditDocument(input.Object.Before),
			After:  auditDocument(input.Object.After),
		},
	}

	_, err := r.collection.InsertOne(ctx, document)
	return err
}

func (r *ChangeRepository) ListRecent(ctx context.Context, limit int64) ([]dto.RecentChangeDto, error) {
	records := make([]dto.RecentChangeDto, 0)
	cursor, err := r.collection.Find(ctx, bson.D{}, options.Find().
		SetSort(bson.D{{Key: "time", Value: -1}, {Key: "_id", Value: -1}}).
		SetLimit(limit).
		SetProjection(bson.M{"time": 1, "user_id": 1, "resource": 1, "resource_id": 1, "action": 1, "object": 1}))
	if err != nil {
		return records, persistenceError(err)
	}
	defer cursor.Close(ctx)
	var changes []entities.Change
	if err := cursor.All(ctx, &changes); err != nil {
		return records, persistenceError(err)
	}
	userNames, err := r.userNames(ctx, changes)
	if err != nil {
		return records, err
	}
	for _, change := range changes {
		records = append(records, dto.RecentChangeDto{
			ID: change.ID.Hex(), Time: change.Time, UserID: change.UserID, UserName: userNames[change.UserID],
			Resource: change.Resource, ResourceID: change.ResourceID, Action: change.Action,
			ChangedFields: changedFields(change.Object.Before, change.Object.After),
		})
	}
	return records, nil
}

func (r *ChangeRepository) userNames(ctx context.Context, changes []entities.Change) (map[string]string, error) {
	names := make(map[string]string)
	ids := make([]bson.ObjectID, 0, len(changes))
	for _, change := range changes {
		if _, seen := names[change.UserID]; seen {
			continue
		}
		names[change.UserID] = ""
		id, err := bson.ObjectIDFromHex(change.UserID)
		if err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return names, nil
	}
	cursor, err := r.db.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": ids}},
		options.Find().SetProjection(bson.M{"name": 1}))
	if err != nil {
		return nil, persistenceError(err)
	}
	defer cursor.Close(ctx)
	var users []struct {
		ID   bson.ObjectID `bson:"_id"`
		Name string        `bson:"name"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, persistenceError(err)
	}
	for _, user := range users {
		names[user.ID.Hex()] = user.Name
	}
	return names, nil
}

func changedFields(before, after bson.M) []dto.ChangedFieldDto {
	fields := make([]dto.ChangedFieldDto, 0)
	keys := make([]string, 0, len(before)+len(after))
	seen := make(map[string]bool)
	for _, document := range []bson.M{before, after} {
		for key := range document {
			if !seen[key] {
				keys = append(keys, key)
				seen[key] = true
			}
		}
	}
	slices.Sort(keys)
	for _, key := range keys {
		if key == "id" || key == "created_at" || key == "deactivated_at" || isSensitiveAuditField(key) {
			continue
		}
		oldValue, newValue := before[key], after[key]
		if reflect.DeepEqual(oldValue, newValue) {
			continue
		}
		oldText, oldSafe := auditScalar(oldValue)
		newText, newSafe := auditScalar(newValue)
		if oldSafe && newSafe {
			fields = append(fields, dto.ChangedFieldDto{Field: key, Before: oldText, After: newText})
		}
	}
	return fields
}

func auditScalar(value any) (string, bool) {
	switch v := value.(type) {
	case nil:
		return "", true
	case string, bool, int, int32, int64, uint, uint32, uint64, float32, float64:
		return fmt.Sprint(v), true
	case time.Time:
		return v.Format(time.RFC3339), true
	default:
		return "", false
	}
}

func auditDocument(value any) bson.M {
	converted := auditValue(reflect.ValueOf(value))
	if document, ok := converted.(bson.M); ok {
		return document
	}
	return bson.M{"value": converted}
}

func auditValue(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}

	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	if value.CanInterface() {
		if timestamp, ok := value.Interface().(time.Time); ok {
			return timestamp
		}
	}

	switch value.Kind() {
	case reflect.Struct:
		typ := value.Type()
		document := bson.M{}
		for i := 0; i < value.NumField(); i++ {
			field := typ.Field(i)
			if field.PkgPath != "" {
				continue
			}

			name := strings.Split(field.Tag.Get("json"), ",")[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = field.Name
			}
			if isSensitiveAuditField(name) {
				document[name] = "[redacted]"
				continue
			}
			document[name] = auditValue(value.Field(i))
		}
		return document

	case reflect.Map:
		document := bson.M{}
		iter := value.MapRange()
		for iter.Next() {
			key := iter.Key()
			if key.Kind() != reflect.String {
				continue
			}
			name := key.String()
			if isSensitiveAuditField(name) {
				document[name] = "[redacted]"
				continue
			}
			document[name] = auditValue(iter.Value())
		}
		return document

	case reflect.Slice, reflect.Array:
		values := make([]any, value.Len())
		for i := 0; i < value.Len(); i++ {
			values[i] = auditValue(value.Index(i))
		}
		return values

	default:
		if value.CanInterface() {
			return value.Interface()
		}
		return nil
	}
}

func isSensitiveAuditField(name string) bool {
	switch strings.ToLower(name) {
	case "password", "password_hash", "token", "access_token", "refresh_token":
		return true
	default:
		return false
	}
}
