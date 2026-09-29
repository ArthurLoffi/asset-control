//go:build integration

package main

import (
	"asset-control/src/main/adapter/outbound/data/database"
	"asset-control/src/main/adapter/outbound/data/repository"
	"asset-control/src/main/domain/dto"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestDeactivationSeparatesActiveRecords(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	if err := database.EnsureSchema(ctx, db); err != nil {
		t.Fatal(err)
	}
	router := testRouter(t, db)
	for _, resource := range []string{"employees", "machines", "phones", "monitors", "printers", "emails"} {
		t.Run(resource, func(t *testing.T) {
			collection := resource
			if resource == "emails" {
				collection = database.UsableEmailsCollection
			}
			createdAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			// Legacy records have no lifecycle field; one also has an explicit null.
			for id := 1; id <= 3; id++ {
				row := bson.M{"id": id, "name": "Same", "user": "Same", "employee_id": 9, "created_at": createdAt, "imported_field": "preserved"}
				if id == 3 {
					row["deactivated_at"] = nil
				}
				if _, err := db.Collection(collection).InsertOne(ctx, row); err != nil {
					t.Fatal(err)
				}
			}
			var active []json.RawMessage
			if err := json.Unmarshal(callAPI(t, router, "GET", "/api/"+resource, nil, 200), &active); err != nil || len(active) != 3 {
				t.Fatalf("legacy records: %s, %v", active, err)
			}
			first := object(t, callAPI(t, router, "POST", "/api/"+resource+"/deactivate?id=1", nil, 200))
			stamp, err := time.Parse(time.RFC3339Nano, first["deactivated_at"].(string))
			if err != nil || time.Since(stamp) > time.Minute {
				t.Fatalf("invalid timestamp: %v", first)
			}
			repeated := object(t, callAPI(t, router, "POST", "/api/"+resource+"/deactivate?id=1", nil, 200))
			if first["deactivated_at"] != repeated["deactivated_at"] {
				t.Fatal("retry changed the timestamp")
			}
			callAPI(t, router, "GET", "/api/"+resource+"/id?id=1", nil, http.StatusNotFound)
			callAPI(t, router, "POST", "/api/"+resource+"/deactivate?id=99", nil, http.StatusNotFound)
			if err := json.Unmarshal(callAPI(t, router, "GET", "/api/"+resource, nil, 200), &active); err != nil || len(active) != 2 {
				t.Fatalf("active list: %s, %v", active, err)
			}
			if resource == "machines" || resource == "phones" || resource == "monitors" {
				if err := json.Unmarshal(callAPI(t, router, "GET", "/api/"+resource+"/employee_id?employee_id=9", nil, 200), &active); err != nil || len(active) != 2 {
					t.Fatalf("assignment list: %s, %v", active, err)
				}
			}
			if resource == "employees" || resource == "emails" {
				query := "name"
				if resource == "emails" {
					query = "user"
				}
				found := object(t, callAPI(t, router, "GET", "/api/"+resource+"/"+query+"?"+query+"=Same", nil, 200))
				if found["id"] != float64(2) {
					t.Fatalf("search selected inactive record: %v", found)
				}
			}
			var stored bson.M
			if err := db.Collection(collection).FindOne(ctx, bson.M{"id": 1}).Decode(&stored); err != nil {
				t.Fatal(err)
			}
			if stored["imported_field"] != "preserved" || stored["created_at"] != bson.NewDateTimeFromTime(createdAt) {
				t.Fatalf("lost original fields: %v", stored)
			}
			callAPI(t, router, "POST", "/api/"+resource+"/deactivate?id=2", nil, 200)
			var page dto.InactivePageDto
			if err := json.Unmarshal(callAPI(t, router, "GET", "/api/"+resource+"/inactive?limit=1", nil, 200), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Records) != 1 || page.Records[0].ID != 1 || page.NextAfter == nil || *page.NextAfter != 1 {
				t.Fatalf("first page: %+v", page)
			}
			if err := json.Unmarshal(callAPI(t, router, "GET", "/api/"+resource+"/inactive?limit=1&after=1", nil, 200), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Records) != 1 || page.Records[0].ID != 2 || page.NextAfter != nil {
				t.Fatalf("last page: %+v", page)
			}
			for _, query := range []string{"limit=0", "limit=101", "after=-1", "limit=no"} {
				callAPI(t, router, "GET", "/api/"+resource+"/inactive?"+query, nil, 400)
			}
		})
	}
	// A normal PUT cannot edit or resurrect an inactive record.
	callAPI(t, router, "PUT", "/api/employees?id=1", map[string]any{"name": "Changed", "email": "a@example.com", "phone": "123", "branch_id": 1, "department_id": 1}, 404)
	callAPI(t, router, "POST", "/api/emails/deactivate?id=2147483648", nil, 400)
}

func TestConcurrentDeactivationPreservesTimestamp(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	if _, err := db.Collection("machines").InsertOne(ctx, bson.M{"id": 1}); err != nil {
		t.Fatal(err)
	}
	repo := repository.NewDeactivationRepository(db)
	stamps := make(chan time.Time, 8)
	var group sync.WaitGroup
	for range 8 {
		group.Go(func() {
			result, err := repo.Deactivate(ctx, "machines", 1)
			if err != nil {
				t.Error(err)
				return
			}
			stamps <- result.DeactivatedAt
		})
	}
	group.Wait()
	close(stamps)
	var first time.Time
	for stamp := range stamps {
		if first.IsZero() {
			first = stamp
		}
		if !first.Equal(stamp) {
			t.Fatal("concurrent retry changed timestamp")
		}
	}
}

func TestEmployeeDeactivationUnassignsActiveEquipment(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	if _, err := db.Collection("employees").InsertOne(ctx, bson.M{"id": 5}); err != nil {
		t.Fatal(err)
	}
	for _, resource := range []string{"machines", "monitors", "phones"} {
		_, err := db.Collection(resource).InsertMany(ctx, []any{
			bson.M{"id": 1, "employee_id": 5},
			bson.M{"id": 2, "employee_id": 6},
			bson.M{"id": 3, "employee_id": 5, "deactivated_at": time.Now().UTC()},
			bson.M{"id": 4, "employee_id": 5, "employee_ids": bson.A{5, 6}},
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	router := testRouter(t, db)
	first := object(t, callAPI(t, router, "POST", "/api/employees/deactivate?id=5", nil, 200))
	repeated := object(t, callAPI(t, router, "POST", "/api/employees/deactivate?id=5", nil, 200))
	if first["deactivated_at"] != repeated["deactivated_at"] {
		t.Fatal("retry changed employee deactivation timestamp")
	}
	for _, resource := range []string{"machines", "monitors", "phones"} {
		for id, want := range map[int]any{1: nil, 2: int32(6), 3: int32(5), 4: int64(6)} {
			var record bson.M
			if err := db.Collection(resource).FindOne(ctx, bson.M{"id": id}).Decode(&record); err != nil {
				t.Fatal(err)
			}
			if record["employee_id"] != want {
				t.Fatalf("%s ID %d assignment: got %v, want %v", resource, id, record["employee_id"], want)
			}
			if id == 4 {
				ids, ok := record["employee_ids"].(bson.A)
				if !ok || len(ids) != 1 || ids[0] != int64(6) {
					t.Fatalf("%s shared assignment lost: %v", resource, record["employee_ids"])
				}
			}
		}
	}
}

func TestUnassignedMachineSearchUsesActiveScope(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	_, err := db.Collection("machines").InsertMany(ctx, []any{
		bson.M{"id": 1, "employee_id": nil},
		bson.M{"id": 2},
		bson.M{"id": 3, "employee_id": 5},
		bson.M{"id": 4, "employee_id": nil, "deactivated_at": time.Now().UTC()},
		bson.M{"id": 5, "employee_id": nil, "employee_ids": bson.A{5}},
		bson.M{"id": 6, "employee_id": nil, "employee_ids": bson.A{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var machines []dto.MachineDto
	if err := json.Unmarshal(callAPI(t, testRouter(t, db), "GET", "/api/machines/unassigned", nil, 200), &machines); err != nil {
		t.Fatal(err)
	}
	if len(machines) != 3 || machines[0].ID != 1 || machines[1].ID != 2 || machines[2].ID != 6 {
		t.Fatalf("unassigned machines: %+v", machines)
	}
}

func TestUnassignedMonitorSearchUsesActiveScope(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	_, err := db.Collection("monitors").InsertMany(ctx, []any{
		bson.M{"id": 1, "employee_id": nil},
		bson.M{"id": 2},
		bson.M{"id": 3, "employee_id": 5},
		bson.M{"id": 4, "employee_id": nil, "deactivated_at": time.Now().UTC()},
		bson.M{"id": 5, "employee_id": nil, "employee_ids": bson.A{5}},
		bson.M{"id": 6, "employee_id": nil, "employee_ids": bson.A{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var monitors []dto.MonitorDto
	if err := json.Unmarshal(callAPI(t, testRouter(t, db), "GET", "/api/monitors/unassigned", nil, 200), &monitors); err != nil {
		t.Fatal(err)
	}
	if len(monitors) != 3 || monitors[0].ID != 1 || monitors[1].ID != 2 || monitors[2].ID != 6 {
		t.Fatalf("unassigned monitors: %+v", monitors)
	}
}
