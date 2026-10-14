package database

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math"
	"time"

	"asset-control/internal/core/dto"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const DefaultPageSize int64 = 20

// Collection stores documents with an application-allocated numeric id, a
// created_at timestamp, and soft delete through deactivated_at.
type Collection[TDocument any] struct {
	collection *mongo.Collection
	counters   *mongo.Collection
	maxID      int64
	notFound   func(id uint64) error
}

// NewCollection reports missing ids with notFound, the entity's domain error.
func NewCollection[TDocument any](db *mongo.Database, name string, notFound func(id uint64) error) Collection[TDocument] {
	return Collection[TDocument]{collection: db.Collection(name), counters: db.Collection(CountersCollection), maxID: math.MaxInt64, notFound: notFound}
}

// WithInt32IDs keeps the imported BSON int32 representation and range for new ids.
func (c Collection[TDocument]) WithInt32IDs() Collection[TDocument] {
	c.maxID = math.MaxInt32
	return c
}

func (c Collection[TDocument]) Database() *mongo.Database { return c.collection.Database() }

func (c Collection[TDocument]) FindPage(ctx context.Context, filter bson.M, sort bson.D, request dto.PaginatedRequestDTO) (dto.PaginatedResponseDTO[TDocument], error) {
	return FindPage[TDocument](ctx, c.collection, activeFilter(filter), sort, request, DefaultPageSize, nil)
}

func (c Collection[TDocument]) FindByID(ctx context.Context, id uint64) (TDocument, error) {
	var document TDocument
	err := c.collection.FindOne(ctx, activeFilter(bson.M{"id": id})).Decode(&document)
	return document, c.persistenceError(err, id)
}

func (c Collection[TDocument]) Create(ctx context.Context, document TDocument) (TDocument, error) {
	var created TDocument
	fields, err := documentFields(document)
	if err != nil {
		return created, err
	}
	var counter struct {
		Seq int64 `bson:"seq"`
	}
	err = c.counters.FindOneAndUpdate(ctx,
		bson.M{"_id": c.collection.Name(), "seq": bson.M{"$lt": c.maxID}},
		bson.M{"$inc": bson.M{"seq": int64(1)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&counter)
	if err != nil {
		return created, fmt.Errorf("allocate %s ID (initialize schema before writes): %w", c.collection.Name(), err)
	}
	fields["id"] = counter.Seq
	if c.maxID == math.MaxInt32 {
		fields["id"] = int32(counter.Seq)
	}
	fields["created_at"] = time.Now().UTC().Truncate(time.Millisecond)
	// Decode before inserting so conversion errors cannot report a failed write.
	encoded, err := bson.Marshal(fields)
	if err != nil {
		return created, err
	}
	if err := bson.Unmarshal(encoded, &created); err != nil {
		return created, err
	}
	if _, err := c.collection.InsertOne(ctx, fields); err != nil {
		return created, PersistenceError(err)
	}
	return created, nil
}

// Update uses $set, so unknown imported fields and created_at are retained and
// explicit nulls/false values are written. It never upserts or reactivates.
func (c Collection[TDocument]) Update(ctx context.Context, id uint64, document TDocument) (TDocument, error) {
	var updated TDocument
	fields, err := documentFields(document)
	if err != nil {
		return updated, err
	}
	err = c.collection.FindOneAndUpdate(ctx, activeFilter(bson.M{"id": id}), bson.M{"$set": fields},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updated)
	return updated, c.persistenceError(err, id)
}

// Deactivate records the database time exactly once with a conditional atomic
// update, even when requests race or a client retries after losing the response.
func (c Collection[TDocument]) Deactivate(ctx context.Context, id uint64) (dto.DeactivationResponseDTO, error) {
	var deactivation dto.DeactivationResponseDTO
	err := c.collection.FindOneAndUpdate(ctx, activeFilter(bson.M{"id": id}),
		bson.M{"$currentDate": bson.M{"deactivated_at": true}},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&deactivation)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = c.collection.FindOne(ctx, bson.M{"id": id, "deactivated_at": bson.M{"$type": "date"}}).Decode(&deactivation)
	}
	return deactivation, c.persistenceError(err, id)
}

// FindInactivePage lists only lifecycle metadata, never the stored document.
func (c Collection[TDocument]) FindInactivePage(ctx context.Context, request dto.PaginatedRequestDTO) (dto.PaginatedResponseDTO[dto.DeactivationResponseDTO], error) {
	return FindPage[dto.DeactivationResponseDTO](ctx, c.collection, bson.M{"deactivated_at": bson.M{"$type": "date"}},
		Sort(request, nil), request, DefaultPageSize, bson.M{"_id": 0, "id": 1, "deactivated_at": 1})
}

func (c Collection[TDocument]) persistenceError(err error, id uint64) error {
	if errors.Is(err, mongo.ErrNoDocuments) && c.notFound != nil {
		return c.notFound(id)
	}
	return PersistenceError(err)
}

// activeFilter copies the filter so callers cannot drop the active scope. Null
// equality also matches imported records without lifecycle metadata.
func activeFilter(filter bson.M) bson.M {
	active := maps.Clone(filter)
	if active == nil {
		active = bson.M{}
	}
	active["deactivated_at"] = nil
	return active
}

// documentFields drops fields that ordinary writes must never change.
func documentFields(document any) (bson.M, error) {
	encoded, err := bson.Marshal(document)
	if err != nil {
		return nil, err
	}
	var fields bson.M
	if err := bson.Unmarshal(encoded, &fields); err != nil {
		return nil, err
	}
	for _, immutable := range []string{"_id", "id", "created_at", "deactivated_at"} {
		delete(fields, immutable)
	}
	return fields, nil
}

// WithinTransaction runs fn in one MongoDB transaction (requires a replica set).
func WithinTransaction(ctx context.Context, db *mongo.Database, fn func(context.Context) error) error {
	session, err := db.Client().StartSession()
	if err != nil {
		return PersistenceError(err)
	}
	defer session.EndSession(context.Background())
	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (any, error) { return nil, fn(txCtx) })
	return err
}
