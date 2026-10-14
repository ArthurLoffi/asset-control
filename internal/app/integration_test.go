//go:build integration

package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"asset-control/internal/core/database"
	coredto "asset-control/internal/core/dto"
	"asset-control/internal/core/exception"
	"asset-control/internal/core/resource"
	"asset-control/internal/modules/inventory/domain/dto"
	inventorydata "asset-control/internal/modules/inventory/infrastructure/outbound/data/adapter"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Every test owns a fresh database; a supplied URI never selects a database to drop.
func testDatabase(t *testing.T) *mongo.Database {
	t.Helper()
	uri := os.Getenv("MONGODB_TEST_URI")
	if uri == "" {
		t.Fatal("MONGODB_TEST_URI is required when running integration tests")
	}
	db, err := database.Connect(context.Background(), database.Config{
		URI: uri, Database: "asset_control_test_" + bson.NewObjectID().Hex(), Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.Drop(ctx); err != nil {
			t.Errorf("drop test database: %v", err)
		}
		if err := database.Disconnect(db.Client()); err != nil {
			t.Errorf("disconnect: %v", err)
		}
	})
	return db
}

func ensureSchema(t *testing.T, db *mongo.Database) {
	t.Helper()
	if err := EnsureSchema(context.Background(), db); err != nil {
		t.Fatal(err)
	}
}

// callAPI returns the raw response body after checking the status.
func callAPI(t *testing.T, router http.Handler, method, target string, input any, status int) json.RawMessage {
	t.Helper()
	var body []byte
	if input != nil {
		var err error
		body, err = json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, target, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != status {
		t.Fatalf("%s %s: got %d, want %d: %s", method, target, recorder.Code, status, recorder.Body.String())
	}
	return recorder.Body.Bytes()
}

func object(t *testing.T, data json.RawMessage) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func page(t *testing.T, data json.RawMessage) ([]map[string]any, coredto.PagingDTO) {
	t.Helper()
	var result struct {
		Content []map[string]any  `json:"content"`
		Paging  coredto.PagingDTO `json:"paging"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result.Content == nil {
		t.Fatalf("content must never be null: %s", data)
	}
	return result.Content, result.Paging
}

func ids(records []map[string]any) []float64 {
	result := make([]float64, 0, len(records))
	for _, record := range records {
		result = append(result, record["id"].(float64))
	}
	return result
}

func TestMongoDBAPIRoundTrips(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	importedAt := time.Date(2025, 3, 4, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		collection   string
		input        map[string]any
		changeField  string
		optionalFlag string
	}{
		{"employees", map[string]any{"name": "Ana [Ops]", "email": "ana@example.com", "phone": "123", "branch_id": 1, "department_id": 2}, "name", ""},
		{"machines", map[string]any{"user": "Ana", "type": "Notebook", "model": "Laptop", "serial_number": "M1", "has_mouse": true, "employee_id": 42}, "model", "has_mouse"},
		{"monitors", map[string]any{"model": "Monitor", "has_monitor": true, "employee_id": 42}, "model", "has_monitor"},
		{"phones", map[string]any{"model": "Phone", "imei": "123", "serial_number": "P1", "has_chip": true, "employee_id": 42}, "model", "has_chip"},
		{"printers", map[string]any{"asset_number": 12, "model": "Printer", "ip": "192.168.100.100", "branch_id": 1, "departments_id": 2}, "model", ""},
	}
	for _, tc := range cases {
		imported := bson.M{"_id": bson.NewObjectID(), "id": int64(42), "created_at": importedAt}
		for key, value := range tc.input {
			imported[key] = value
		}
		if _, err := db.Collection(tc.collection).InsertOne(ctx, imported); err != nil {
			t.Fatal(err)
		}
	}
	for name, id := range map[string]int64{"branches": 1, "departments": 2} {
		if _, err := db.Collection(name).InsertOne(ctx, bson.M{"id": id, "name": name + " name"}); err != nil {
			t.Fatal(err)
		}
	}
	ensureSchema(t, db)
	gin.SetMode(gin.TestMode)
	router := testRouter(t, db)
	for _, tc := range cases {
		t.Run(tc.collection, func(t *testing.T) {
			path := "/api/" + tc.collection
			imported := object(t, callAPI(t, router, "GET", path+"/42", nil, 200))
			if imported["id"] != float64(42) || imported["created_at"] != importedAt.Format(time.RFC3339) {
				t.Fatalf("imported document not preserved: %v", imported)
			}
			created := object(t, callAPI(t, router, "POST", path, tc.input, 201))
			if created["id"] != float64(43) || created["created_at"] == nil {
				t.Fatalf("bad generated fields: %v", created)
			}
			if tc.optionalFlag != "" && created["situation"] != "Em uso" {
				t.Fatalf("SQL default missing: %v", created)
			}
			if _, exposed := created["_id"]; exposed {
				t.Fatal("MongoDB _id exposed through the API")
			}
			saved := object(t, callAPI(t, router, "GET", path+"/43", nil, 200))
			if saved[tc.changeField] != tc.input[tc.changeField] {
				t.Fatalf("wrong persisted value: %v", saved)
			}
			records, paging := page(t, callAPI(t, router, "GET", path, nil, 200))
			if len(records) != 2 || records[0]["id"] != float64(42) || paging.TotalElements != 2 || paging.Size != 20 {
				t.Fatalf("list=%v paging=%+v", records, paging)
			}
			_, err := db.Collection(tc.collection).UpdateOne(ctx, bson.M{"id": 43}, bson.M{"$set": bson.M{"migration_note": "keep me"}})
			if err != nil {
				t.Fatal(err)
			}
			tc.input[tc.changeField] = "Updated"
			tc.input["id"] = 999 // Client cannot change generated identifiers.
			tc.input["created_at"] = "2000-01-01T00:00:00Z"
			if tc.optionalFlag != "" {
				tc.input[tc.optionalFlag] = false
				tc.input["employee_id"] = nil
			}
			updated := object(t, callAPI(t, router, "PUT", path+"/43", tc.input, 200))
			if updated["id"] != float64(43) || updated["created_at"] != created["created_at"] || updated[tc.changeField] != "Updated" {
				t.Fatalf("incorrect update: %v", updated)
			}
			var raw bson.M
			if err := db.Collection(tc.collection).FindOne(ctx, bson.M{"id": 43}).Decode(&raw); err != nil {
				t.Fatal(err)
			}
			if raw["migration_note"] != "keep me" {
				t.Fatal("update removed unknown imported fields")
			}
			if tc.optionalFlag != "" {
				if value, exists := raw["employee_id"]; !exists || value != nil || raw[tc.optionalFlag] != false {
					t.Fatalf("null/false update lost: %v", raw)
				}
				related, _ := page(t, callAPI(t, router, "GET", path+"?employee_id=42", nil, 200))
				if len(related) != 1 || related[0]["id"] != float64(42) {
					t.Fatalf("relationship filter=%v", related)
				}
				if empty, _ := page(t, callAPI(t, router, "GET", path+"?employee_id=999", nil, 200)); len(empty) != 0 {
					t.Fatalf("no match must be []: %v", empty)
				}
			}
			notFound := object(t, callAPI(t, router, "GET", path+"/999", nil, 404))
			if notFound["code"] != strings.ToUpper(strings.TrimSuffix(tc.collection, "s"))+"_001" {
				t.Fatalf("entity not-found code: %v", notFound)
			}
			callAPI(t, router, "PUT", path+"/999", tc.input, 404)
			callAPI(t, router, "POST", path, map[string]any{}, 400)
			callAPI(t, router, "PUT", path+"/43", map[string]any{}, 400)
			callAPI(t, router, "GET", path+"/0", nil, 400)
			count, err := db.Collection(tc.collection).CountDocuments(ctx, bson.M{})
			if err != nil || count != 2 {
				t.Fatalf("invalid request or missing update wrote a document: %d, %v", count, err)
			}
		})
	}
	for path, want := range map[string]string{"/api/branches/1": "branches name", "/api/departments/2": "departments name"} {
		if data := object(t, callAPI(t, router, "GET", path, nil, 200)); data["name"] != want || data["id"] == nil {
			t.Fatalf("lookup=%v", data)
		}
	}
	if branches, _ := page(t, callAPI(t, router, "GET", "/api/branches?name=BRANCH", nil, 200)); len(branches) != 1 {
		t.Fatalf("branch name filter: %v", branches)
	}
	callAPI(t, router, "GET", "/api/branches/999", nil, 404)
	callAPI(t, router, "GET", "/api/departments/999", nil, 404)
	found, _ := page(t, callAPI(t, router, "GET", "/api/employees?name="+url.QueryEscape("ana [ops]"), nil, 200))
	if len(found) != 1 || found[0]["id"] != float64(42) {
		t.Fatalf("literal case-insensitive name search: %v", found)
	}
	if none, _ := page(t, callAPI(t, router, "GET", "/api/employees?name="+url.QueryEscape(".*"), nil, 200)); len(none) != 0 {
		t.Fatalf("regex characters must be literal: %v", none)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	_, err := inventorydata.NewEmployeeDataAdapter(db).FindAll(canceled, dto.EmployeeRequestDTO{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
}

func TestListingContract(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	names := []string{"Carla", "ana", "Bruno", "Ana Maria", "Diego"}
	for index, name := range names {
		if _, err := db.Collection("employees").InsertOne(ctx, bson.M{"id": index + 1, "name": name, "branch_id": index%2 + 1}); err != nil {
			t.Fatal(err)
		}
	}
	router := testRouter(t, db)
	for _, tc := range []struct {
		query  string
		ids    []float64
		paging coredto.PagingDTO
	}{
		{"", []float64{1, 2, 3, 4, 5}, coredto.PagingDTO{TotalElements: 5, Size: 20, TotalPages: 1, HasContent: true}},
		{"?page_size=2", []float64{1, 2}, coredto.PagingDTO{TotalElements: 5, Size: 2, TotalPages: 3, HasContent: true, HasNext: true}},
		{"?page_size=2&page=2", []float64{5}, coredto.PagingDTO{TotalElements: 5, Size: 2, Page: 2, TotalPages: 3, HasContent: true, HasPrevious: true}},
		{"?page_size=2&page=9", []float64{}, coredto.PagingDTO{TotalElements: 5, Size: 2, Page: 9, TotalPages: 3, HasPrevious: true}},
		{"?direction=desc", []float64{5, 4, 3, 2, 1}, coredto.PagingDTO{TotalElements: 5, Size: 20, TotalPages: 1, HasContent: true}},
		{"?sort_by=name", []float64{2, 4, 3, 1, 5}, coredto.PagingDTO{TotalElements: 5, Size: 20, TotalPages: 1, HasContent: true}},
		{"?name=ANA", []float64{2, 4}, coredto.PagingDTO{TotalElements: 2, Size: 20, TotalPages: 1, HasContent: true}},
		{"?branch_id=1&direction=DESC", []float64{5, 3, 1}, coredto.PagingDTO{TotalElements: 3, Size: 20, TotalPages: 1, HasContent: true}},
	} {
		records, paging := page(t, callAPI(t, router, "GET", "/api/employees"+tc.query, nil, 200))
		if fmt.Sprint(ids(records)) != fmt.Sprint(tc.ids) || paging != tc.paging {
			t.Errorf("%s: ids=%v paging=%+v", tc.query, ids(records), paging)
		}
	}
	for _, query := range []string{"page=-1", "page_size=101", "direction=up", "sort_by=password", "page=x", "branch_id=x"} {
		body := object(t, callAPI(t, router, "GET", "/api/employees?"+query, nil, 400))
		if body["code"] != exception.ValidationCode {
			t.Errorf("%s: %v", query, body)
		}
	}
}

func TestMongoDBConcurrentIDsAndRestart(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	if _, err := db.Collection("employees").InsertOne(ctx, bson.M{"id": int64(100), "name": "Imported"}); err != nil {
		t.Fatal(err)
	}
	ensureSchema(t, db)
	create := resource.NewUseCases("employee", inventorydata.NewEmployeeDataAdapter(db), nil).Create
	const workers = 32
	createdIDs := make(chan uint64, workers)
	failures := make(chan error, workers+4)
	var wg sync.WaitGroup
	for i := range workers {
		wg.Go(func() {
			created, err := create.Execute(ctx, dto.EmployeeInputDTO{Name: fmt.Sprintf("Employee %d", i), Email: fmt.Sprintf("employee%d@example.com", i), BranchID: 1, DepartmentID: 2})
			if err != nil {
				failures <- err
				return
			}
			createdIDs <- created.ID
		})
	}
	// Simulate additional application instances starting while writes are active.
	for range 4 {
		wg.Go(func() {
			if err := EnsureSchema(ctx, db); err != nil {
				failures <- err
			}
		})
	}
	wg.Wait()
	close(createdIDs)
	close(failures)
	for err := range failures {
		t.Error(err)
	}
	seen := map[uint64]bool{}
	for id := range createdIDs {
		if id <= 100 || seen[id] {
			t.Errorf("duplicate/import-conflicting ID: %d", id)
		}
		seen[id] = true
	}
	if len(seen) != workers {
		t.Fatalf("got %d IDs, want %d", len(seen), workers)
	}
	if _, err := db.Collection("employees").DeleteOne(ctx, bson.M{"id": 100 + workers}); err != nil {
		t.Fatal(err)
	}
	ensureSchema(t, db)
	created, err := create.Execute(ctx, dto.EmployeeInputDTO{Name: "After restart", Email: "restart@example.com", BranchID: 1, DepartmentID: 2})
	if err != nil || created.ID != 101+workers {
		t.Fatalf("counter moved backwards: %+v, %v", created, err)
	}
	// Force an imported-ID collision to verify the adapter maps duplicate keys.
	if _, err := db.Collection(database.CountersCollection).UpdateOne(ctx, bson.M{"_id": "employees"}, bson.M{"$set": bson.M{"seq": int64(99)}}); err != nil {
		t.Fatal(err)
	}
	_, err = create.Execute(ctx, dto.EmployeeInputDTO{Name: "Collision", Email: "collision@example.com", BranchID: 1, DepartmentID: 2})
	if !errors.Is(err, exception.ErrConflict) {
		t.Fatalf("duplicate error not mapped: %v", err)
	}
}

func TestMongoDBEmptyCollections(t *testing.T) {
	db := testDatabase(t)
	ensureSchema(t, db)
	gin.SetMode(gin.TestMode)
	router := testRouter(t, db)
	for _, name := range []string{"employees", "machines", "monitors", "phones", "printers", "emails", "branches", "departments", "changes"} {
		want := `{"content":[],"paging":{"total_elements":0,"size":20,"page":0,"total_pages":0,"has_content":false,"has_next":false,"has_previous":false}}`
		if result := callAPI(t, router, "GET", "/api/"+name, nil, 200); string(result) != want {
			t.Fatalf("%s: %s", name, result)
		}
	}
}
