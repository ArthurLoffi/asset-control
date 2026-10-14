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
	"asset-control/internal/core/exception"
	"asset-control/internal/modules/inventory/application/usecase"
	"asset-control/internal/modules/inventory/domain/dto"
	repository "asset-control/internal/modules/inventory/infrastructure/outbound/data/adapter"
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
	var response struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	return response.Data
}

func object(t *testing.T, data json.RawMessage) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
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
	if err := database.EnsureSchema(ctx, db); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	router := testRouter(t, db)
	for _, tc := range cases {
		t.Run(tc.collection, func(t *testing.T) {
			path := "/api/" + tc.collection
			imported := object(t, callAPI(t, router, "GET", path+"/id?id=42", nil, 200))
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
				t.Fatal("MongoDB _id exposed through legacy API")
			}
			saved := object(t, callAPI(t, router, "GET", path+"/id?id=43", nil, 200))
			if saved[tc.changeField] != tc.input[tc.changeField] {
				t.Fatalf("wrong persisted value: %v", saved)
			}
			list := callAPI(t, router, "GET", path, nil, 200)
			var records []map[string]any
			if err := json.Unmarshal(list, &records); err != nil || len(records) != 2 {
				t.Fatalf("list=%s err=%v", list, err)
			}
			if records[0]["id"] != float64(42) {
				t.Fatal("list order not deterministic")
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
			updated := object(t, callAPI(t, router, "PUT", path+"?id=43", tc.input, 200))
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
				related := callAPI(t, router, "GET", path+"/employee_id?employee_id=42", nil, 200)
				if err := json.Unmarshal(related, &records); err != nil || len(records) != 1 || records[0]["id"] != float64(42) {
					t.Fatalf("relationship filter=%s err=%v", related, err)
				}
				empty := callAPI(t, router, "GET", path+"/employee_id?employee_id=999", nil, 200)
				if string(empty) != "[]" {
					t.Fatalf("empty collection must be [], got %s", empty)
				}
			}
			callAPI(t, router, "GET", path+"/id?id=999", nil, 404)
			callAPI(t, router, "PUT", path+"?id=999", tc.input, 404)
			callAPI(t, router, "POST", path, map[string]any{}, 400)
			callAPI(t, router, "PUT", path+"?id=43", map[string]any{}, 400)
			callAPI(t, router, "GET", path+"/id?id=0", nil, 400)
			count, err := db.Collection(tc.collection).CountDocuments(ctx, bson.M{})
			if err != nil || count != 2 {
				t.Fatalf("invalid request or missing update wrote a document: %d, %v", count, err)
			}
		})
	}
	for path, want := range map[string]string{"/api/utils/branch?id=1": "branches name", "/api/utils/department?id=2": "departments name"} {
		data := object(t, callAPI(t, router, "GET", path, nil, 200))
		if data["name"] != want {
			t.Fatalf("lookup=%v", data)
		}
	}
	callAPI(t, router, "GET", "/api/utils/branch?id=999", nil, 404)
	callAPI(t, router, "GET", "/api/utils/department?id=999", nil, 404)
	found := object(t, callAPI(t, router, "GET", "/api/employees/name?name="+url.QueryEscape("ana [ops]"), nil, 200))
	if found["id"] != float64(42) {
		t.Fatalf("literal case-insensitive name search: %v", found)
	}
	callAPI(t, router, "GET", "/api/employees/name?name="+url.QueryEscape(".*"), nil, 404)
	callAPI(t, router, "GET", "/api/employees/name", nil, 400)
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	_, err := repository.NewEmployeeRepository(db).GetAllEmployees(canceled)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
}

func TestMongoDBConcurrentIDsAndRestart(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	if _, err := db.Collection("employees").InsertOne(ctx, bson.M{"id": int64(100), "name": "Imported"}); err != nil {
		t.Fatal(err)
	}
	if err := database.EnsureSchema(ctx, db); err != nil {
		t.Fatal(err)
	}
	uc := usecase.NewEmployeeUseCase(repository.NewEmployeeRepository(db))
	const workers = 32
	ids := make(chan uint64, workers)
	failures := make(chan error, workers+4)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			created, err := uc.CreateEmployee(ctx, dto.CreateEmployeeDto{Name: fmt.Sprintf("Employee %d", i), Email: fmt.Sprintf("employee%d@example.com", i), BranchID: 1, DepartmentID: 2})
			if err != nil {
				failures <- err
				return
			}
			ids <- created.ID
		}(i)
	}
	// Simulate additional application instances starting while writes are active.
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := database.EnsureSchema(ctx, db); err != nil {
				failures <- err
			}
		}()
	}
	wg.Wait()
	close(ids)
	close(failures)
	for err := range failures {
		t.Error(err)
	}
	seen := map[uint64]bool{}
	for id := range ids {
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
	if err := database.EnsureSchema(ctx, db); err != nil {
		t.Fatal(err)
	}
	created, err := uc.CreateEmployee(ctx, dto.CreateEmployeeDto{Name: "After restart", Email: "restart@example.com", BranchID: 1, DepartmentID: 2})
	if err != nil || created.ID != 101+workers {
		t.Fatalf("counter moved backwards: %+v, %v", created, err)
	}
	// Force an imported-ID collision to verify the API repository maps duplicate keys.
	if _, err := db.Collection(database.CountersCollection).UpdateOne(ctx, bson.M{"_id": "employees"}, bson.M{"$set": bson.M{"seq": int64(99)}}); err != nil {
		t.Fatal(err)
	}
	_, err = uc.CreateEmployee(ctx, dto.CreateEmployeeDto{Name: "Collision", Email: "collision@example.com", BranchID: 1, DepartmentID: 2})
	if !errors.Is(err, exception.ErrConflict) {
		t.Fatalf("duplicate error not mapped: %v", err)
	}
}

func TestMongoDBEmptyCollections(t *testing.T) {
	db := testDatabase(t)
	if err := database.EnsureSchema(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	router := testRouter(t, db)
	for _, name := range []string{"employees", "machines", "monitors", "phones", "printers"} {
		result := callAPI(t, router, "GET", "/api/"+name, nil, 200)
		if string(result) != "[]" {
			t.Fatalf("%s: %s", name, result)
		}
	}
}
