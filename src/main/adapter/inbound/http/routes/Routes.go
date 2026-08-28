package routes

import (
	"asset-control/src/main/adapter/inbound/http/controller"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, employeeController *controller.EmployeeController) {
	api := r.Group("/api")

	e := api.Group("/employee")
	{
		e.POST("", employeeController.CreateEmployeeController())
		e.GET("/id", employeeController.FindByIDController())
		e.GET("/name", employeeController.FindByNameController())
		e.GET("", employeeController.GetAllEmployees())
		e.PUT("", employeeController.UpdateEmployeeController())
	}
}
