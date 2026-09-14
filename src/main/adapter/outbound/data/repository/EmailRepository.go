package repository

import (
	"context"
	"math"
	"regexp"

	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EmailRepository struct {
	documents documentRepository[entities.Email]
}

var _ ports.EmailsPorts = (*EmailRepository)(nil)

func NewEmailRepository(db *mongo.Database) *EmailRepository {
	return &EmailRepository{documents: newDocumentRepository[entities.Email](db, (entities.Email{}).CollectionName())}
}

func emailDocument(input dto.CreateEmailDto) entities.Email {
	return entities.Email{User: input.User, Email: input.Email, PasswordHash: input.PasswordHash}
}

func (r *EmailRepository) GetAllEmails(ctx context.Context) ([]dto.EmailDto, error) {
	documents, err := r.documents.find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.EmailDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.EmailDto(document))
	}
	return result, nil
}

func (r *EmailRepository) CreateEmail(ctx context.Context, input dto.CreateEmailDto) (dto.EmailDto, error) {
	document, err := r.documents.createWithIDLimit(ctx, emailDocument(input), math.MaxInt32)
	return dto.EmailDto(document), err
}

func (r *EmailRepository) UpdateEmail(ctx context.Context, id int32, input dto.CreateEmailDto) (dto.EmailDto, error) {
	document, err := r.documents.update(ctx, uint64(id), emailDocument(input))
	return dto.EmailDto(document), err
}

func (r *EmailRepository) FindByID(ctx context.Context, id int32) (dto.EmailDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"id": id})
	return dto.EmailDto(document), err
}

func (r *EmailRepository) FindByUser(ctx context.Context, user string) (dto.EmailDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"user": bson.Regex{Pattern: regexp.QuoteMeta(user), Options: "i"}})
	return dto.EmailDto(document), err
}
