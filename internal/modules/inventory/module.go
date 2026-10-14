package inventory

import (
	"context"

	"asset-control/internal/core/audit"
	"asset-control/internal/core/database"
	"asset-control/internal/core/resource"
	inventoryapi "asset-control/internal/modules/inventory/infrastructure/inbound/http/adapter"
	inventorydata "asset-control/internal/modules/inventory/infrastructure/outbound/data/adapter"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Register wires each entity's data adapter, use cases, and HTTP adapter.
// Resource names are the singular values stored in audit records.
func Register(router gin.IRouter, db *mongo.Database, recorder audit.Recorder) {
	inventoryapi.NewEmployeeAPIAdapter(resource.NewUseCases("employee", inventorydata.NewEmployeeDataAdapter(db), recorder)).RegisterRoutes(router)
	inventoryapi.NewMachineAPIAdapter(resource.NewUseCases("machine", inventorydata.NewMachineDataAdapter(db), recorder)).RegisterRoutes(router)
	inventoryapi.NewMonitorAPIAdapter(resource.NewUseCases("monitor", inventorydata.NewMonitorDataAdapter(db), recorder)).RegisterRoutes(router)
	inventoryapi.NewPhoneAPIAdapter(resource.NewUseCases("phone", inventorydata.NewPhoneDataAdapter(db), recorder)).RegisterRoutes(router)
	inventoryapi.NewPrinterAPIAdapter(resource.NewUseCases("printer", inventorydata.NewPrinterDataAdapter(db), recorder)).RegisterRoutes(router)
	inventoryapi.NewEmailAPIAdapter(resource.NewUseCases("email", inventorydata.NewEmailDataAdapter(db), recorder)).RegisterRoutes(router)
}

// EnsureSchema creates indexes and seeds the id counters; run before accepting writes.
func EnsureSchema(ctx context.Context, db *mongo.Database) error {
	assignmentIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "employee_id", Value: 1}, {Key: "deactivated_at", Value: 1}, {Key: "id", Value: 1}}},
		{Keys: bson.D{{Key: "employee_ids", Value: 1}, {Key: "deactivated_at", Value: 1}, {Key: "id", Value: 1}}},
	}
	for _, name := range []string{schema.EmployeeCollection, schema.PrinterCollection, schema.EmailCollection} {
		if err := database.EnsureResource(ctx, db, name); err != nil {
			return err
		}
	}
	for _, name := range []string{schema.MachineCollection, schema.MonitorCollection, schema.PhoneCollection} {
		if err := database.EnsureResource(ctx, db, name, assignmentIndexes...); err != nil {
			return err
		}
	}
	return nil
}
