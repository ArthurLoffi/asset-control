package adapter

import (
	"context"

	"asset-control/internal/core/database"
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// assignedCollections hold equipment that references employees.
var assignedCollections = []string{schema.MachineCollection, schema.MonitorCollection, schema.PhoneCollection}

// assignmentFilter matches the legacy employee_id and the employee_ids list.
func assignmentFilter(filter dto.AssignmentFilterDTO) bson.M {
	conditions := bson.A{}
	if len(filter.EmployeeIDs) > 0 {
		anyOf := bson.M{"$in": filter.EmployeeIDs}
		conditions = append(conditions, bson.M{"$or": bson.A{bson.M{"employee_id": anyOf}, bson.M{"employee_ids": anyOf}}})
	}
	if filter.IsUnassigned {
		conditions = append(conditions, bson.M{
			"employee_id": nil,
			"$or":         bson.A{bson.M{"employee_ids": nil}, bson.M{"employee_ids": bson.M{"$size": 0}}},
		})
	}
	if len(conditions) == 0 {
		return bson.M{}
	}
	return bson.M{"$and": conditions}
}

// unassignEmployee removes only the deactivated employee from active equipment;
// the equipment itself stays active.
func unassignEmployee(db *mongo.Database) func(context.Context, uint64) error {
	return func(ctx context.Context, employeeID uint64) error {
		for _, name := range assignedCollections {
			if err := removeAssignment(ctx, db.Collection(name), employeeID); err != nil {
				return err
			}
		}
		return nil
	}
}

func removeAssignment(ctx context.Context, collection *mongo.Collection, employeeID uint64) error {
	filter := assignmentFilter(dto.AssignmentFilterDTO{EmployeeIDs: []uint64{employeeID}})
	filter["deactivated_at"] = nil
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"id": 1, "employee_id": 1, "employee_ids": 1}))
	if err != nil {
		return database.PersistenceError(err)
	}
	var records []struct {
		ID          uint64   `bson:"id"`
		EmployeeID  *uint64  `bson:"employee_id"`
		EmployeeIDs []uint64 `bson:"employee_ids"`
	}
	if err := cursor.All(ctx, &records); err != nil {
		return database.PersistenceError(err)
	}
	for _, record := range records {
		remaining := withoutEmployeeID(record.EmployeeID, record.EmployeeIDs, employeeID)
		var first *uint64
		if len(remaining) > 0 {
			first = &remaining[0]
		}
		_, err := collection.UpdateOne(ctx, bson.M{"id": record.ID, "deactivated_at": nil},
			bson.M{"$set": bson.M{"employee_id": first, "employee_ids": remaining}})
		if err != nil {
			return database.PersistenceError(err)
		}
	}
	return nil
}

func withoutEmployeeID(single *uint64, many []uint64, removed uint64) []uint64 {
	if many == nil && single != nil {
		many = []uint64{*single}
	}
	remaining := make([]uint64, 0, len(many))
	for _, id := range many {
		if id != removed {
			remaining = append(remaining, id)
		}
	}
	return remaining
}
