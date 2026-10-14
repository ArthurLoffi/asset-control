package database

import (
	"errors"

	"asset-control/internal/core/exception"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func PersistenceError(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return exception.ErrNotFound
	case mongo.IsDuplicateKeyError(err):
		return exception.ErrConflict
	default:
		return err
	}
}
