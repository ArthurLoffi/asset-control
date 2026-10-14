package audit

import (
	"context"

	coreaudit "asset-control/internal/core/audit"
	"asset-control/internal/modules/audit/application/usecase"
	auditapi "asset-control/internal/modules/audit/infrastructure/inbound/http/adapter"
	auditdata "asset-control/internal/modules/audit/infrastructure/outbound/data/adapter"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Register mounts the change feed and returns the recorder other modules use
// to write audit records.
func Register(router gin.IRouter, db *mongo.Database) coreaudit.Recorder {
	changeDataAdapter := auditdata.NewChangeDataAdapter(db)
	auditapi.NewChangeAPIAdapter(usecase.NewFindAllChangesUseCase(changeDataAdapter)).RegisterRoutes(router)
	return changeDataAdapter
}

func EnsureSchema(ctx context.Context, db *mongo.Database) error {
	return auditdata.NewChangeDataAdapter(db).EnsureSchema(ctx)
}
