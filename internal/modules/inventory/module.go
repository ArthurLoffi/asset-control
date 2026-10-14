package inventory

import (
	"asset-control/internal/core/audit"
	"asset-control/internal/modules/inventory/application/usecase"
	inventoryapi "asset-control/internal/modules/inventory/infrastructure/inbound/http/adapter"
	inventorydata "asset-control/internal/modules/inventory/infrastructure/outbound/data/adapter"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Register(router gin.IRouter, db *mongo.Database, recorder audit.Recorder) {
	employeeController := inventoryapi.NewEmployeeController(usecase.NewEmployeeUseCase(inventorydata.NewEmployeeRepository(db), recorder))
	machineController := inventoryapi.NewMachineController(usecase.NewMachineUseCase(inventorydata.NewMachineRepository(db), recorder))
	phoneController := inventoryapi.NewPhoneController(usecase.NewPhoneUseCase(inventorydata.NewPhoneRepository(db), recorder))
	monitorController := inventoryapi.NewMonitorController(usecase.NewMonitorUseCase(inventorydata.NewMonitorRepository(db), recorder))
	printerController := inventoryapi.NewPrinterController(usecase.NewPrinterUseCase(inventorydata.NewPrinterRepository(db), recorder))
	emailController := inventoryapi.NewEmailController(usecase.NewEmailUseCase(inventorydata.NewEmailRepository(db), recorder))
	deactivationController := inventoryapi.NewDeactivationController(usecase.NewDeactivationUseCase(inventorydata.NewDeactivationRepository(db)))

	for _, resource := range []string{"employees", "machines", "phones", "monitors", "printers", "emails"} {
		router.POST("/"+resource+"/deactivate", deactivationController.Deactivate(resource))
		router.GET("/"+resource+"/inactive", deactivationController.ListInactive(resource))
	}

	emails := router.Group("/emails")
	emails.GET("", emailController.GetAllEmailsController())
	emails.POST("", emailController.CreateEmailController())
	emails.PUT("", emailController.UpdateEmailController())
	emails.GET("/id", emailController.FindByIDController())
	emails.GET("/user", emailController.FindByUserController())

	employees := router.Group("/employees")
	employees.POST("", employeeController.CreateEmployeeController())
	employees.GET("/id", employeeController.FindByIDController())
	employees.GET("/name", employeeController.FindByNameController())
	employees.GET("", employeeController.GetAllEmployeesController())
	employees.PUT("", employeeController.UpdateEmployeeController())

	machines := router.Group("/machines")
	machines.POST("", machineController.CreateMachineController())
	machines.GET("/id", machineController.FindByIDController())
	machines.GET("/employee_id", machineController.FindByEmployeeIDController())
	machines.GET("/unassigned", machineController.GetUnassignedMachinesController())
	machines.GET("", machineController.GetAllMachinesController())
	machines.PUT("", machineController.UpdateMachineController())

	phones := router.Group("/phones")
	phones.POST("", phoneController.CreatePhoneController())
	phones.GET("/id", phoneController.FindByIDController())
	phones.GET("/employee_id", phoneController.FindByEmployeeIDController())
	phones.GET("", phoneController.GetAllPhonesController())
	phones.PUT("", phoneController.UpdatePhoneController())

	monitors := router.Group("/monitors")
	monitors.POST("", monitorController.CreateMonitorController())
	monitors.GET("/id", monitorController.FindByIDController())
	monitors.GET("/employee_id", monitorController.FindByEmployeeIDController())
	monitors.GET("/unassigned", monitorController.GetUnassignedMonitorsController())
	monitors.GET("", monitorController.GetAllMonitorsController())
	monitors.PUT("", monitorController.UpdateMonitorController())

	printers := router.Group("/printers")
	printers.POST("", printerController.CreatePrinterController())
	printers.GET("/id", printerController.FindByIDController())
	printers.GET("", printerController.GetAllPrintersController())
	printers.PUT("", printerController.UpdatePrinterController())
}
