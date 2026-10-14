package audit

import (
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
	changeDataAdapter := auditdata.NewChangeRepository(db)
	changeAPIAdapter := auditapi.NewRecentChangeController(usecase.NewRecentChangeUseCase(changeDataAdapter))
	router.GET("/changes/recent", changeAPIAdapter.ListRecent)
	return changeDataAdapter
}
