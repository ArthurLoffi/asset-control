package repository

import "go.mongodb.org/mongo-driver/v2/bson"

func effectiveEmployeeIDs(single *uint64, many []uint64) []uint64 {
	if many != nil {
		return many
	}
	if single != nil {
		return []uint64{*single}
	}
	return []uint64{}
}

func withoutEmployeeID(ids []uint64, removed uint64) []uint64 {
	remaining := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id != removed {
			remaining = append(remaining, id)
		}
	}
	return remaining
}

func assignedTo(id uint64) bson.M {
	return bson.M{"$or": bson.A{bson.M{"employee_id": id}, bson.M{"employee_ids": id}}}
}

func unassignedFilter() bson.M {
	return bson.M{
		"employee_id": nil,
		"$or":         bson.A{bson.M{"employee_ids": nil}, bson.M{"employee_ids": bson.M{"$size": 0}}},
	}
}
