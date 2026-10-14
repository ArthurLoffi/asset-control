package organization

import (
	"context"

	"asset-control/internal/core/database"
	"asset-control/internal/core/resource"
	organizationapi "asset-control/internal/modules/organization/infrastructure/inbound/http/adapter"
	organizationdata "asset-control/internal/modules/organization/infrastructure/outbound/data/adapter"
	"asset-control/internal/modules/organization/infrastructure/outbound/data/schema"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Register(router gin.IRouter, db *mongo.Database) {
	organizationapi.NewBranchAPIAdapter(resource.NewReadUseCases("branch", organizationdata.NewBranchDataAdapter(db))).RegisterRoutes(router)
	organizationapi.NewDepartmentAPIAdapter(resource.NewReadUseCases("department", organizationdata.NewDepartmentDataAdapter(db))).RegisterRoutes(router)
}

func EnsureSchema(ctx context.Context, db *mongo.Database) error {
	for _, name := range []string{schema.BranchCollection, schema.DepartmentCollection} {
		if err := database.EnsureLookup(ctx, db, name); err != nil {
			return err
		}
	}
	return nil
}
