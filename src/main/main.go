package main

import (
	"asset-control/src/main/adapter/inbound/http/controller"
	"asset-control/src/main/adapter/inbound/http/routes"
	_ "asset-control/src/main/adapter/inbound/http/swagger"
	"asset-control/src/main/adapter/outbound/data/database"
	"asset-control/src/main/adapter/outbound/data/repository"
	"asset-control/src/main/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

// @title           Asset control API
// @version         v0.2.0
// @description     Está é uma API utilizada para controlar os equipamentos de TI.

// @contact.name   Suporte BNTG
// @contact.email  suporte@brasilnovo.com

// @host      localhost:8080
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	db, err := database.NewMySQLConnection()
	if err != nil {
		log.Fatal(err)
	}

	employeeRepository := repository.NewEmployeeRepository(db)
	employeeUseCase := usecase.NewEmployeeUseCase(employeeRepository)
	employeeController := controller.NewEmployeeController(employeeUseCase)

	machineRepository := repository.NewMachineRepository(db)
	machineUseCase := usecase.NewMachineUseCase(machineRepository)
	machineController := controller.NewMachineController(machineUseCase)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.Routes(r, employeeController, machineController)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
