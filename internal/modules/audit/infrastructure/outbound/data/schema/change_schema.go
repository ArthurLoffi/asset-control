package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChangeObject struct {
	Before bson.M `bson:"before"`
	After  bson.M `bson:"after"`
}

type Change struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	Time       time.Time     `bson:"time"`
	UserID     string        `bson:"user_id"`
	Resource   string        `bson:"resource"`
	ResourceID int64         `bson:"resource_id"`
	Action     string        `bson:"action"`
	Object     ChangeObject  `bson:"object"`
}

func (Change) CollectionName() string { return "changes" }
