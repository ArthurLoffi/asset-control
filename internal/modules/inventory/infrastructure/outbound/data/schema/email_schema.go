package schema

import (
	"time"

	"asset-control/internal/core/database"
)

const EmailCollection = database.UsableEmailsCollection

// EmailSortColumns maps the public sort_by value to the stored field.
var EmailSortColumns = map[string]string{"id": "id", "user": "user", "created_at": "created_at"}

type EmailDocument struct {
	ID           int32     `bson:"id"`
	User         string    `bson:"user"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
}
