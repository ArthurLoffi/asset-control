package controller

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/usecase"
	"net/http"
	"strconv"

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

func (ctl *EmployeeController) UpdateEmployeeController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'id' is required"})
			return
		}

		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}

		var employee dto.EmployeeDto
		if err := c.ShouldBindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := ctl.employeeUseCase.UpdateEmployee(id, employee)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "employee updated successfully",
			"data":   updated,
		})
	}
}

func (ctl *EmployeeController) GetAllEmployees() gin.HandlerFunc {
	return func(c *gin.Context) {
		employees, err := ctl.employeeUseCase.GetAllEmployees()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "employees retrieved successfully",
			"data":   employees,
		})
	}
}

func (ctl *EmployeeController) FindByNameController() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "name is required",
			})
		}

		employee, err := ctl.employeeUseCase.FindByName(name)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "user found successfully",
			"data":   employee,
		})
	}
}

func (ctl *EmployeeController) FindByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'id' is required"})
			return
		}

		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}

		employee, err := ctl.employeeUseCase.FindByID(id)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "user found successfully",
			"data":   employee,
		})
	}
}
