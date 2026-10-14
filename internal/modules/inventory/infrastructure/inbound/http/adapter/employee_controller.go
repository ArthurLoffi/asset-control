package adapter

import (
	"asset-control/internal/core/exception"
	"asset-control/internal/core/httpx"
	"net/http"
	"strings"

	"asset-control/internal/modules/inventory/application/usecase"
	"asset-control/internal/modules/inventory/domain/dto"

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

// CreateEmployee godoc
// @Summary      Cria um novo funcionário
// @Description  Cria um novo funcionário a partir dos dados enviados no corpo da requisição.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        employee  body      dto.CreateEmployeeDto  true  "Dados do funcionário"
// @Success      201       {object}  map[string]interface{}  "Funcionário criado com sucesso"
// @Failure      400       {object}  exception.ErrorResponse  "Dados inválidos"
// @Failure      500       {object}  exception.ErrorResponse  "Erro interno do servidor"
// @Router       /employees [post]
func (ctl *EmployeeController) CreateEmployeeController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee dto.CreateEmployeeDto

		if err := httpx.BindJSON(c, &employee); err != nil {
			httpx.WriteError(c, err)
			return
		}

		created, err := ctl.employeeUseCase.CreateEmployee(c.Request.Context(), employee)
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

// UpdateEmployee godoc
// @Summary      Atualiza um funcionário
// @Description  Atualiza os dados de um funcionário existente através do ID.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id        query     uint64           true  "ID do funcionário"
// @Param        employee  body      dto.CreateEmployeeDto  true  "Dados atualizados do funcionário"
// @Success      200       {object}  map[string]interface{}  "Funcionário atualizado com sucesso"
// @Failure      400       {object}  exception.ErrorResponse  "ID inválido ou dados inválidos"
// @Failure      404       {object}  exception.ErrorResponse  "Funcionário não encontrado"
// @Router       /employees [put]
func (ctl *EmployeeController) UpdateEmployeeController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := httpx.QueryID(c, "id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		var employee dto.CreateEmployeeDto
		if err := httpx.BindJSON(c, &employee); err != nil {
			httpx.WriteError(c, err)
			return
		}

		updated, err := ctl.employeeUseCase.UpdateEmployee(c.Request.Context(), id, employee)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "employee updated successfully",
			"data":   updated,
		})
	}
}

// GetAllEmployees godoc
// @Summary      Lista todos os funcionários
// @Description  Retorna todos os funcionários cadastrados.
// @Tags         employees
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Funcionários encontrados com sucesso"
// @Failure      500  {object}  exception.ErrorResponse  "Erro interno do servidor"
// @Router       /employees [get]
func (ctl *EmployeeController) GetAllEmployeesController() gin.HandlerFunc {
	return func(c *gin.Context) {
		employees, err := ctl.employeeUseCase.GetAllEmployees(c.Request.Context())
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "employees retrieved successfully",
			"data":   employees,
		})
	}
}

// FindByName godoc
// @Summary      Busca funcionário por nome
// @Description  Busca um funcionário pelo nome informado como query parameter.
// @Tags         employees
// @Produce      json
// @Param        name  query  string  true  "Nome do funcionário"
// @Success      200  {object}  map[string]interface{}  "Funcionário encontrado com sucesso"
// @Failure      400  {object}  exception.ErrorResponse  "Nome não informado"
// @Failure      404  {object}  exception.ErrorResponse  "Funcionário não encontrado"
// @Router       /employees/name [get]
func (ctl *EmployeeController) FindByNameController() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.TrimSpace(c.Query("name"))

		if name == "" {
			httpx.WriteError(c, exception.NewValidationError("name is required"))
			return
		}

		employee, err := ctl.employeeUseCase.FindByName(c.Request.Context(), name)

		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "user found successfully",
			"data":   employee,
		})
	}
}

// FindByID godoc
// @Summary      Busca funcionário por ID
// @Description  Busca um funcionário pelo ID informado como query parameter.
// @Tags         employees
// @Produce      json
// @Param        id  query  uint64  true  "ID do funcionário"
// @Success      200  {object}  map[string]interface{}  "Funcionário encontrado com sucesso"
// @Failure      400  {object}  exception.ErrorResponse  "ID inválido ou não informado"
// @Failure      404  {object}  exception.ErrorResponse  "Funcionário não encontrado"
// @Router       /employees/id [get]
func (ctl *EmployeeController) FindByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := httpx.QueryID(c, "id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		employee, err := ctl.employeeUseCase.FindByID(c.Request.Context(), id)

		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "user found successfully",
			"data":   employee,
		})
	}
}
