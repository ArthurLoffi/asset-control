package entities

import (
	"time"

	"asset-control/src/main/infrastructure/outbound/data/database"
)

type Email struct {
	ID           int32     `bson:"id" json:"id"`
	User         string    `bson:"user" json:"user"`
	Email        string    `bson:"email" json:"email"`
	PasswordHash string    `bson:"password_hash" json:"password_hash"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

func (Email) CollectionName() string {
	return database.UsableEmailsCollection
}
