package adapter

import (
	"asset-control/internal/core/httpx"
	"net/http"

	"asset-control/internal/modules/inventory/application/usecase"
	"asset-control/internal/modules/inventory/domain/dto"

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

// CreateMachine godoc
// @Summary      Cria uma nova máquina
// @Description  Cria uma nova máquina a partir dos dados enviados no corpo da requisição.
// @Tags         machines
// @Accept       json
// @Produce      json
// @Param        machine  body      dto.CreateMachineDto  true  "Dados da máquina"
// @Success      201      {object}  map[string]interface{}  "Máquina criada com sucesso"
// @Failure      400      {object}  map[string]string  "Dados inválidos"
// @Failure      500      {object}  map[string]string  "Erro interno do servidor"
// @Router       /machines [post]
func (ctl *MachineController) CreateMachineController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var machine dto.CreateMachineDto

		if err := c.ShouldBindJSON(&machine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		created, err := ctl.machineUseCase.CreateMachine(c.Request.Context(), machine)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "created",
			"data":   created,
		})
	}
}

// UpdateMachine godoc
// @Summary      Atualiza uma máquina
// @Description  Atualiza os dados de uma máquina existente através do ID informado como query parameter.
// @Tags         machines
// @Accept       json
// @Produce      json
// @Param        id       query     uint64          true  "ID da máquina"
// @Param        machine  body      dto.CreateMachineDto  true  "Dados atualizados da máquina"
// @Success      200      {object}  map[string]interface{}  "Máquina atualizada com sucesso"
// @Failure      400      {object}  map[string]string  "ID ou dados inválidos"
// @Failure      404      {object}  map[string]string  "Máquina não encontrada"
// @Router       /machines [put]
func (ctl *MachineController) UpdateMachineController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'id' is required"})
			return
		}

		id, err := httpx.ParseID(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}

		var machine dto.CreateMachineDto
		if err := c.ShouldBindJSON(&machine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := ctl.machineUseCase.UpdateMachine(c.Request.Context(), id, machine)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machine updated successfully",
			"data":   updated,
		})
	}
}

// GetAllMachines godoc
// @Summary      Lista todas as máquinas
// @Description  Retorna todas as máquinas cadastradas.
// @Tags         machines
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Máquinas encontradas com sucesso"
// @Failure      500  {object}  map[string]string  "Erro interno do servidor"
// @Router       /machines [get]
func (ctl *MachineController) GetAllMachinesController() gin.HandlerFunc {
	return func(c *gin.Context) {
		machines, err := ctl.machineUseCase.GetAllMachines(c.Request.Context())
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machines retrieved successfully",
			"data":   machines,
		})
	}
}

// FindMachineByID godoc
// @Summary      Busca máquina por ID
// @Description  Busca uma máquina pelo ID informado como query parameter.
// @Tags         machines
// @Produce      json
// @Param        id  query  uint64  true  "ID da máquina"
// @Success      200  {object}  map[string]interface{}  "Máquina encontrada com sucesso"
// @Failure      400  {object}  map[string]string  "ID inválido ou não informado"
// @Failure      404  {object}  map[string]string  "Máquina não encontrada"
// @Router       /machines/id [get]
func (ctl *MachineController) FindByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'id' is required"})
			return
		}

		id, err := httpx.ParseID(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}

		machine, err := ctl.machineUseCase.FindByID(c.Request.Context(), id)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machine found successfully",
			"data":   machine,
		})
	}
}

// FindMachineByEmployeeID godoc
// @Summary      Busca máquinas por funcionário
// @Description  Retorna as máquinas vinculadas ao funcionário informado pelo ID.
// @Tags         machines
// @Produce      json
// @Param        employee_id  query  uint64  true  "ID do funcionário"
// @Success      200  {object}  map[string]interface{}  "Máquinas encontradas com sucesso"
// @Failure      400  {object}  map[string]string  "ID do funcionário inválido ou não informado"
// @Failure      404  {object}  map[string]string  "Nenhuma máquina encontrada para o funcionário"
// @Router       /machines/employee_id [get]
func (ctl *MachineController) FindByEmployeeIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("employee_id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'employee_id' is required"})
			return
		}

		employeeId, err := httpx.ParseID(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'employee_id' value"})
			return
		}

		machine, err := ctl.machineUseCase.FindByEmployeeID(c.Request.Context(), employeeId)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "machine found successfully",
			"data":   machine,
		})
	}
}

// GetUnassignedMachines godoc
// @Summary      List active machines without an assigned employee
// @Tags         machines
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /machines/unassigned [get]
func (ctl *MachineController) GetUnassignedMachinesController() gin.HandlerFunc {
	return func(c *gin.Context) {
		machines, err := ctl.machineUseCase.GetUnassignedMachines(c.Request.Context())
		if err != nil {
			httpx.WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "unassigned machines retrieved successfully",
			"data":   machines,
		})
	}
}
