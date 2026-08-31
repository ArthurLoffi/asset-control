package controller

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MachineController struct {
	machineUseCase *usecase.MachineUseCase
}

func NewMachineController(machineUseCase *usecase.MachineUseCase) *MachineController {
	return &MachineController{
		machineUseCase: machineUseCase,
	}
}

func (ctl *MachineController) CreateMachineController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var machine dto.MachineDto

		if err := c.ShouldBindJSON(&machine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := ctl.machineUseCase.CreateMachine(machine); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "created",
			"data":   machine,
		})
	}
}

func (ctl *MachineController) UpdateMachineController() gin.HandlerFunc {
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

		var machine dto.MachineDto
		if err := c.ShouldBindJSON(&machine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := ctl.machineUseCase.UpdateMachine(id, machine)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machine updated successfully",
			"data":   updated,
		})
	}
}

func (ctl *MachineController) GetAllMachinesController() gin.HandlerFunc {
	return func(c *gin.Context) {
		machines, err := ctl.machineUseCase.GetAllMachines()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machines retrieved successfully",
			"data":   machines,
		})
	}
}

func (ctl *MachineController) FindByIDController() gin.HandlerFunc {
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

		machine, err := ctl.machineUseCase.FindByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machine found successfully",
			"data":   machine,
		})
	}
}

func (ctl *MachineController) FindByEmployeeIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("employee_id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'employee_id' is required"})
			return
		}

		employeeId, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'employee_id' value"})
			return
		}

		machine, err := ctl.machineUseCase.FindByEmployeeID(employeeId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machine found successfully",
			"data":   machine,
		})
	}
}
