package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	ChangeCollection = "changes"
	// UserCollection is owned by asset-control-auth; only names are read.
	UserCollection = "users"
)

type ChangeObjectDocument struct {
	Before bson.M `bson:"before"`
	After  bson.M `bson:"after"`
}

type ChangeDocument struct {
	ID         bson.ObjectID        `bson:"_id,omitempty"`
	Time       time.Time            `bson:"time"`
	UserID     string               `bson:"user_id"`
	Resource   string               `bson:"resource"`
	ResourceID int64                `bson:"resource_id"`
	Action     string               `bson:"action"`
	Object     ChangeObjectDocument `bson:"object"`
}
