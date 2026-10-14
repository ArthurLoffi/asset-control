package adapter

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
	"time"
)

func TestActiveFilterIncludesLegacyRecordsAndPreservesCriteria(t *testing.T) {
	filter := bson.M{"employee_id": uint64(8)}
	active := activeFilter(filter)
	if active["employee_id"] != uint64(8) {
		t.Fatal("lost caller criteria")
	}
	value, ok := active["deactivated_at"]
	if !ok || value != nil {
		t.Fatal("active reads must match missing or null timestamps")
	}
	if _, exists := filter["deactivated_at"]; exists {
		t.Fatal("modified caller filter")
	}
}

func TestNormalWritesCannotChangeLifecycleMetadata(t *testing.T) {
	fields, err := documentFields(bson.M{"id": 1, "model": "original", "deactivated_at": time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := fields["deactivated_at"]; exists {
		t.Fatal("normal writes can change lifecycle metadata")
	}
}
