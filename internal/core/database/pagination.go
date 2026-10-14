package database

import (
	"context"
	"regexp"

	"asset-control/internal/core/dto"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

// sortCollation orders text case-insensitively in Portuguese ("ana" next to
// "Ana Maria"). Count uses it too, so both queries match the same documents.
var sortCollation = &options.Collation{Locale: "pt", Strength: 2}

// FindPage counts and finds in parallel. A page past the end is empty, never an error.
func FindPage[TDocument any](
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
	sort bson.D,
	request dto.PaginatedRequestDTO,
	defaultPageSize int64,
	projection bson.M,
) (dto.PaginatedResponseDTO[TDocument], error) {
	pageSize := request.EffectivePageSize(defaultPageSize)
	var totalElements int64
	var documents []TDocument

	queryGroup, queryCtx := errgroup.WithContext(ctx)
	queryGroup.Go(func() error {
		documentCount, err := collection.CountDocuments(queryCtx, filter, options.Count().SetCollation(sortCollation))
		totalElements = documentCount
		return err
	})
	queryGroup.Go(func() error {
		findOptions := options.Find().SetSort(sort).SetSkip(request.Page * pageSize).SetLimit(pageSize).SetCollation(sortCollation)
		if projection != nil {
			findOptions.SetProjection(projection)
		}
		cursor, err := collection.Find(queryCtx, filter, findOptions)
		if err != nil {
			return err
		}
		return cursor.All(queryCtx, &documents)
	})
	if err := queryGroup.Wait(); err != nil {
		return dto.PaginatedResponseDTO[TDocument]{}, PersistenceError(err)
	}
	return dto.NewPaginatedResponseDTO(documents, request.Page, pageSize, totalElements), nil
}

// Sort maps the public sort_by through the entity's allow list and always adds
// id as the tie-breaker, so the order is stable across pages. Client text never
// reaches the database.
func Sort(request dto.PaginatedRequestDTO, sortColumns map[string]string) bson.D {
	order := 1
	if request.Direction == dto.DirectionDesc {
		order = -1
	}
	sortColumn, isSortable := sortColumns[request.SortBy]
	if !isSortable || sortColumn == "id" {
		return bson.D{{Key: "id", Value: order}}
	}
	return bson.D{{Key: sortColumn, Value: order}, {Key: "id", Value: order}}
}

// Contains matches partially and case-insensitively; the text is escaped, so
// regular expression characters are literal.
func Contains(text string) bson.M {
	return bson.M{"$regex": regexp.QuoteMeta(text), "$options": "i"}
}
