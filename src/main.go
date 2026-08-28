package main

import (
	"asset-control/src/main/adapter/inbound/http/controller"
	"asset-control/src/main/adapter/inbound/http/routes"
	"asset-control/src/main/adapter/outbound/data/database"
	"asset-control/src/main/adapter/outbound/data/repository"
	"asset-control/src/main/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewMySQLConnection()
	if err != nil {
		log.Fatal(err)
	}

	employeeRepository := repository.NewEmployeeRepository(db)
	employeeUseCase := usecase.NewEmployeeUseCase(employeeRepository)
	employeeController := controller.NewEmployeeController(employeeUseCase)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.Routes(r, employeeController)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
