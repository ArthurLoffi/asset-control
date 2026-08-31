package routes

import (
	"asset-control/src/main/adapter/inbound/http/controller"

	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.Engine,
	employeeController *controller.EmployeeController,
	machineController *controller.MachineController,
) {
	api := r.Group("/api")

	e := api.Group("/employees")
	{
		e.POST("", employeeController.CreateEmployeeController())
		e.GET("/id", employeeController.FindByIDController())
		e.GET("/name", employeeController.FindByNameController())
		e.GET("", employeeController.GetAllEmployees())
		e.PUT("", employeeController.UpdateEmployeeController())
	}

	m := api.Group("/machines")
	{
		m.POST("", machineController.CreateMachineController())
		m.GET("/id", machineController.FindByIDController())
		m.GET("/employee_id", machineController.FindByEmployeeIDController())
		m.GET("", machineController.GetAllMachinesController())
		m.PUT("", machineController.UpdateMachineController())
	}
}
