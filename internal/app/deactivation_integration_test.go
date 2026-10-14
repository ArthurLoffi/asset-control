//go:build integration

package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"asset-control/internal/core/database"
	inventorydata "asset-control/internal/modules/inventory/infrastructure/outbound/data/adapter"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestDeactivationSeparatesActiveRecords(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	ensureSchema(t, db)
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
			path := "/api/" + resource
			if active, _ := page(t, callAPI(t, router, "GET", path, nil, 200)); len(active) != 3 {
				t.Fatalf("legacy records: %v", active)
			}
			first := object(t, callAPI(t, router, "DELETE", path+"/1", nil, 200))
			stamp, err := time.Parse(time.RFC3339Nano, first["deactivated_at"].(string))
			if err != nil || time.Since(stamp) > time.Minute || first["id"] != float64(1) {
				t.Fatalf("invalid deactivation: %v", first)
			}
			repeated := object(t, callAPI(t, router, "DELETE", path+"/1", nil, 200))
			if first["deactivated_at"] != repeated["deactivated_at"] {
				t.Fatal("retry changed the timestamp")
			}
			callAPI(t, router, "GET", path+"/1", nil, http.StatusNotFound)
			callAPI(t, router, "DELETE", path+"/99", nil, http.StatusNotFound)
			if active, _ := page(t, callAPI(t, router, "GET", path, nil, 200)); fmt.Sprint(ids(active)) != "[2 3]" {
				t.Fatalf("active list: %v", active)
			}
			if resource == "machines" || resource == "phones" || resource == "monitors" {
				if active, _ := page(t, callAPI(t, router, "GET", path+"?employee_id=9", nil, 200)); len(active) != 2 {
					t.Fatalf("assignment list: %v", active)
				}
			}
			if resource == "employees" || resource == "emails" {
				query := map[string]string{"employees": "name", "emails": "user"}[resource]
				if found, _ := page(t, callAPI(t, router, "GET", path+"?"+query+"=Same", nil, 200)); fmt.Sprint(ids(found)) != "[2 3]" {
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
			callAPI(t, router, "DELETE", path+"/2", nil, 200)
			inactive, paging := page(t, callAPI(t, router, "GET", path+"/inactive?page_size=1", nil, 200))
			if len(inactive) != 1 || inactive[0]["id"] != float64(1) || !paging.HasNext || len(inactive[0]) != 2 {
				t.Fatalf("first page: %v %+v", inactive, paging)
			}
			inactive, paging = page(t, callAPI(t, router, "GET", path+"/inactive?page_size=1&page=1", nil, 200))
			if len(inactive) != 1 || inactive[0]["id"] != float64(2) || paging.HasNext {
				t.Fatalf("last page: %v %+v", inactive, paging)
			}
			for _, query := range []string{"page_size=101", "page=-1", "page=no", "sort_by=name"} {
				callAPI(t, router, "GET", path+"/inactive?"+query, nil, 400)
			}
		})
	}
	// A normal PUT cannot edit or resurrect an inactive record.
	callAPI(t, router, "PUT", "/api/employees/1", map[string]any{"name": "Changed", "email": "a@example.com", "phone": "123", "branch_id": 1, "department_id": 1}, 404)
	callAPI(t, router, "DELETE", "/api/emails/2147483648", nil, 404)
}

func TestConcurrentDeactivationPreservesTimestamp(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	if _, err := db.Collection("machines").InsertOne(ctx, bson.M{"id": 1}); err != nil {
		t.Fatal(err)
	}
	machines := inventorydata.NewMachineDataAdapter(db)
	stamps := make(chan time.Time, 8)
	var group sync.WaitGroup
	for range 8 {
		group.Go(func() {
			result, err := machines.Deactivate(ctx, 1)
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
	first := object(t, callAPI(t, router, "DELETE", "/api/employees/5", nil, 200))
	repeated := object(t, callAPI(t, router, "DELETE", "/api/employees/5", nil, 200))
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
				employeeIDs, ok := record["employee_ids"].(bson.A)
				if !ok || len(employeeIDs) != 1 || employeeIDs[0] != int64(6) {
					t.Fatalf("%s shared assignment lost: %v", resource, record["employee_ids"])
				}
			}
		}
	}
}

func TestUnassignedFilterUsesActiveScope(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	router := testRouter(t, db)
	for _, resource := range []string{"machines", "monitors", "phones"} {
		_, err := db.Collection(resource).InsertMany(ctx, []any{
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
		if unassigned, _ := page(t, callAPI(t, router, "GET", "/api/"+resource+"?unassigned=true", nil, 200)); fmt.Sprint(ids(unassigned)) != "[1 2 6]" {
			t.Fatalf("unassigned %s: %v", resource, unassigned)
		}
		if assigned, _ := page(t, callAPI(t, router, "GET", "/api/"+resource+"?employee_id=5&employee_id=7", nil, 200)); fmt.Sprint(ids(assigned)) != "[3 5]" {
			t.Fatalf("assigned %s: %v", resource, assigned)
		}
	}
}
