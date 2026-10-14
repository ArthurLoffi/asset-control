package adapter

import (
	"asset-control/internal/core/httpx"
	"net/http"

	"asset-control/internal/modules/inventory/application/usecase"
	"asset-control/internal/modules/inventory/domain/dto"

	"github.com/gin-gonic/gin"
)

type MonitorController struct {
	monitorUseCase *usecase.MonitorUseCase
}

func NewMonitorController(monitorUseCase *usecase.MonitorUseCase) *MonitorController {
	return &MonitorController{
		monitorUseCase: monitorUseCase,
	}
}

// CreateMonitor godoc
// @Summary      Cria um novo monitor
// @Description  Cria um novo monitor a partir dos dados enviados no corpo da requisição.
// @Tags         monitors
// @Accept       json
// @Produce      json
// @Param        monitor  body      dto.CreateMonitorDto  true  "Dados do monitor"
// @Success      201       {object}  map[string]interface{}  "Monitor criado com sucesso"
// @Failure      400       {object}  exception.ErrorResponse  "Dados inválidos"
// @Failure      500       {object}  exception.ErrorResponse  "Erro interno do servidor"
// @Router       /monitors [post]
func (ctl *MonitorController) CreateMonitorController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var monitor dto.CreateMonitorDto

		if err := httpx.BindJSON(c, &monitor); err != nil {
			httpx.WriteError(c, err)
			return
		}

		created, err := ctl.monitorUseCase.CreateMonitor(c.Request.Context(), monitor)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "monitor created successfully",
			"data":   created,
		})
	}
}

// UpdateMonitor godoc
// @Summary      Atualiza um monitor
// @Description  Atualiza os dados de um monitor existente através do ID.
// @Tags         monitors
// @Accept       json
// @Produce      json
// @Param        id       query     uint64                 true  "ID do monitor"
// @Param        monitor  body      dto.CreateMonitorDto   true  "Dados atualizados do monitor"
// @Success      200      {object}  map[string]interface{}  "Monitor atualizado com sucesso"
// @Failure      400      {object}  exception.ErrorResponse  "ID inválido ou dados inválidos"
// @Failure      404      {object}  exception.ErrorResponse  "Monitor não encontrado"
// @Router       /monitors [put]
func (ctl *MonitorController) UpdateMonitorController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := httpx.QueryID(c, "id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		var monitor dto.CreateMonitorDto

		if err := httpx.BindJSON(c, &monitor); err != nil {
			httpx.WriteError(c, err)
			return
		}

		updated, err := ctl.monitorUseCase.UpdateMonitor(c.Request.Context(), id, monitor)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "monitor updated successfully",
			"data":   updated,
		})
	}
}

// GetAllMonitors godoc
// @Summary      Lista todos os monitores
// @Description  Retorna todos os monitores cadastrados.
// @Tags         monitors
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Monitores encontrados com sucesso"
// @Failure      500  {object}  exception.ErrorResponse  "Erro interno do servidor"
// @Router       /monitors [get]
func (ctl *MonitorController) GetAllMonitorsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		monitors, err := ctl.monitorUseCase.GetAllMonitors(c.Request.Context())

		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "monitors retrieved successfully",
			"data":   monitors,
		})
	}
}

// FindByID godoc
// @Summary      Busca monitor por ID
// @Description  Busca um monitor pelo ID informado como query parameter.
// @Tags         monitors
// @Produce      json
// @Param        id  query  uint64  true  "ID do monitor"
// @Success      200  {object}  map[string]interface{}  "Monitor encontrado com sucesso"
// @Failure      400  {object}  exception.ErrorResponse  "ID inválido ou não informado"
// @Failure      404  {object}  exception.ErrorResponse  "Monitor não encontrado"
// @Router       /monitors/id [get]
func (ctl *MonitorController) FindByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := httpx.QueryID(c, "id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		monitor, err := ctl.monitorUseCase.FindByID(c.Request.Context(), id)

		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "monitor found successfully",
			"data":   monitor,
		})
	}
}

// FindByEmployeeID godoc
// @Summary      Busca monitores por funcionário
// @Description  Retorna todos os monitores associados ao ID de um funcionário.
// @Tags         monitors
// @Produce      json
// @Param        employee_id  query  uint64  true  "ID do funcionário"
// @Success      200  {object}  map[string]interface{}  "Monitores encontrados com sucesso"
// @Failure      400  {object}  exception.ErrorResponse  "ID do funcionário inválido ou não informado"
// @Failure      404  {object}  exception.ErrorResponse  "Nenhum monitor encontrado para o funcionário"
// @Router       /monitors/employee_id [get]
func (ctl *MonitorController) FindByEmployeeIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID, err := httpx.QueryID(c, "employee_id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		monitors, err := ctl.monitorUseCase.FindByEmployeeID(c.Request.Context(), employeeID)

		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "monitors found successfully",
			"data":   monitors,
		})
	}
}

// GetUnassignedMonitors godoc
// @Summary      List active monitors without an assigned employee
// @Tags         monitors
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /monitors/unassigned [get]
func (ctl *MonitorController) GetUnassignedMonitorsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		monitors, err := ctl.monitorUseCase.GetUnassignedMonitors(c.Request.Context())
		if err != nil {
			httpx.WriteError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "unassigned monitors retrieved successfully",
			"data":   monitors,
		})
	}
}
