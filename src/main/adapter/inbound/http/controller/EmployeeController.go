package controller

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	employeeUseCase *usecase.EmployeeUseCase
}

func NewEmployeeController(employeeUseCase *usecase.EmployeeUseCase) *EmployeeController {
	return &EmployeeController{
		employeeUseCase: employeeUseCase,
	}
}

func (ctl *EmployeeController) CreateEmployeeController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee dto.EmployeeDto

		if err := c.ShouldBindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := ctl.employeeUseCase.CreateEmployee(employee); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "created",
			"data":   employee,
		})
	}
}
