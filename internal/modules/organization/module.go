package organization

import (
	"asset-control/internal/modules/organization/application/usecase"
	organizationapi "asset-control/internal/modules/organization/infrastructure/inbound/http/adapter"
	organizationdata "asset-control/internal/modules/organization/infrastructure/outbound/data/adapter"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Register(router gin.IRouter, db *mongo.Database) {
	utilsController := organizationapi.NewUtilsController(usecase.NewUtilsUseCase(organizationdata.NewUtilsRepository(db)))

	utils := router.Group("/utils")
	utils.GET("/branches", utilsController.GetAllBranchs())
	utils.GET("/branch", utilsController.BranchIdToDescription())
	utils.GET("/departments", utilsController.GetAllDepartments())
	utils.GET("/department", utilsController.DepartmentIdToDescription())
}
