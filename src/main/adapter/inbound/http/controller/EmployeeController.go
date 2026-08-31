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

// CreateEmployee godoc
// @Summary      Cria um novo funcionário
// @Description  Cria um novo funcionário a partir dos dados enviados no corpo da requisição.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        employee  body      dto.CreateEmployeeDto  true  "Dados do funcionário"
// @Success      201       {object}  map[string]interface{}  "Funcionário criado com sucesso"
// @Failure      400       {object}  map[string]string  "Dados inválidos"
// @Failure      500       {object}  map[string]string  "Erro interno do servidor"
// @Router       /employees [post]
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

// UpdateEmployee godoc
// @Summary      Atualiza um funcionário
// @Description  Atualiza os dados de um funcionário existente através do ID.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id        query     uint64           true  "ID do funcionário"
// @Param        employee  body      dto.CreateEmployeeDto  true  "Dados atualizados do funcionário"
// @Success      200       {object}  map[string]interface{}  "Funcionário atualizado com sucesso"
// @Failure      400       {object}  map[string]string  "ID inválido ou dados inválidos"
// @Failure      404       {object}  map[string]string  "Funcionário não encontrado"
// @Router       /employees [put]
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

// GetAllEmployees godoc
// @Summary      Lista todos os funcionários
// @Description  Retorna todos os funcionários cadastrados.
// @Tags         employees
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Funcionários encontrados com sucesso"
// @Failure      500  {object}  map[string]string  "Erro interno do servidor"
// @Router       /employees [get]
func (ctl *EmployeeController) GetAllEmployeesController() gin.HandlerFunc {
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

// FindByName godoc
// @Summary      Busca funcionário por nome
// @Description  Busca um funcionário pelo nome informado como query parameter.
// @Tags         employees
// @Produce      json
// @Param        name  query  string  true  "Nome do funcionário"
// @Success      200  {object}  map[string]interface{}  "Funcionário encontrado com sucesso"
// @Failure      400  {object}  map[string]string  "Nome não informado"
// @Failure      404  {object}  map[string]string  "Funcionário não encontrado"
// @Router       /employees/name [get]
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

// FindByID godoc
// @Summary      Busca funcionário por ID
// @Description  Busca um funcionário pelo ID informado como query parameter.
// @Tags         employees
// @Produce      json
// @Param        id  query  uint64  true  "ID do funcionário"
// @Success      200  {object}  map[string]interface{}  "Funcionário encontrado com sucesso"
// @Failure      400  {object}  map[string]string  "ID inválido ou não informado"
// @Failure      404  {object}  map[string]string  "Funcionário não encontrado"
// @Router       /employees/id [get]
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
