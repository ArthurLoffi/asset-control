package routes

import (
	"asset-control/src/main/adapter/inbound/http/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(
	r *gin.Engine,
	employeeController *controller.EmployeeController,
	machineController *controller.MachineController,
	phoneController *controller.PhoneController,
	monitorController *controller.MonitorController,
	printerController *controller.PrinterController,
	utilsController *controller.UtilsController,
	emailController *controller.EmailController,
) {
	api := r.Group("/api")

	{
		e := api.Group("/emails")
		{
			e.GET("", emailController.GetAllEmailsController())
			e.POST("", emailController.CreateEmailController())
			e.PUT("", emailController.UpdateEmailController())
			e.GET("/id", emailController.FindByIDController())
			e.GET("/user", emailController.FindByUserController())
		}
	}

	{
		e := api.Group("/employees")
		{
			e.POST("", employeeController.CreateEmployeeController())
			e.GET("/id", employeeController.FindByIDController())
			e.GET("/name", employeeController.FindByNameController())
			e.GET("", employeeController.GetAllEmployeesController())
			e.PUT("", employeeController.UpdateEmployeeController())
		}
	}

	{
		m := api.Group("/machines")
		{
			m.POST("", machineController.CreateMachineController())
			m.GET("/id", machineController.FindByIDController())
			m.GET("/employee_id", machineController.FindByEmployeeIDController())
			m.GET("", machineController.GetAllMachinesController())
			m.PUT("", machineController.UpdateMachineController())
		}
	}

	{
		p := api.Group("/phones")
		{
			p.POST("", phoneController.CreatePhoneController())
			p.GET("/id", phoneController.FindByIDController())
			p.GET("/employee_id", phoneController.FindByEmployeeIDController())
			p.GET("", phoneController.GetAllPhonesController())
			p.PUT("", phoneController.UpdatePhoneController())
		}
	}

	{
		mo := api.Group("/monitors")
		{
			mo.POST("", monitorController.CreateMonitorController())
			mo.GET("/id", monitorController.FindByIDController())
			mo.GET("/employee_id", monitorController.FindByEmployeeIDController())
			mo.GET("", monitorController.GetAllMonitorsController())
			mo.PUT("", monitorController.UpdateMonitorController())
		}
	}

	{
		pr := api.Group("/printers")
		{
			pr.POST("", printerController.CreatePrinterController())
			pr.GET("/id", printerController.FindByIDController())
			pr.GET("", printerController.GetAllPrintersController())
			pr.PUT("", printerController.UpdatePrinterController())
		}
	}

	{
		utils := api.Group("/utils")
		{
			utils.GET("/branches", utilsController.GetAllBranchs())
			utils.GET("/branch", utilsController.BranchIdToDescription())
			utils.GET("/departments", utilsController.GetAllDepartments())
			utils.GET("/department", utilsController.DepartmentIdToDescription())
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
