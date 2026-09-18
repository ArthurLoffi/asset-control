package repository

import (
	"context"
	"reflect"
	"strings"
	"time"

	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
