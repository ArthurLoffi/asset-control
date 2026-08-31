package controller

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhoneController struct {
	phoneUseCase *usecase.PhoneUseCase
}

func NewPhoneController(phoneUseCase *usecase.PhoneUseCase) *PhoneController {
	return &PhoneController{
		phoneUseCase: phoneUseCase,
	}
}

// CreatePhone godoc
// @Summary Cria um novo telefone
// @Description Cria um novo telefone a partir dos dados enviados no corpo da requisição.
// @Tags phones
// @Accept json
// @Produce json
// @Param phone body dto.CreatePhoneDto true "Dados do telefone"
// @Success 201 {object} map[string]interface{} "Telefone criado com sucesso"
// @Failure 400 {object} map[string]string "Dados inválidos"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /phones [post]
func (ctl *PhoneController) CreatePhoneController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var phone dto.CreatePhoneDto

		if err := c.ShouldBindJSON(&phone); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := ctl.phoneUseCase.CreatePhone(phone); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "created",
			"data":   phone,
		})
	}
}

// UpdatePhone godoc
// @Summary Atualiza um telefone
// @Description Atualiza os dados de um telefone existente através do ID informado como query parameter.
// @Tags phones
// @Accept json
// @Produce json
// @Param id query uint64 true "ID do telefone"
// @Param phone body dto.CreatePhoneDto true "Dados atualizados do telefone"
// @Success 200 {object} map[string]interface{} "Telefone atualizado com sucesso"
// @Failure 400 {object} map[string]string "ID ou dados inválidos"
// @Failure 404 {object} map[string]string "Telefone não encontrado"
// @Router /phones [put]
func (ctl *PhoneController) UpdatePhoneController() gin.HandlerFunc {
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

		var phone dto.CreatePhoneDto
		if err := c.ShouldBindJSON(&phone); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := ctl.phoneUseCase.UpdatePhone(id, phone)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "phone updated successfully",
			"data":   updated,
		})
	}
}

// GetAllPhones godoc
// @Summary Lista todos os telefones
// @Description Retorna todos os telefones cadastrados.
// @Tags phones
// @Produce json
// @Success 200 {object} map[string]interface{} "Telefones encontrados com sucesso"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /phones [get]
func (ctl *PhoneController) GetAllPhonesController() gin.HandlerFunc {
	return func(c *gin.Context) {
		phones, err := ctl.phoneUseCase.GetAllPhones()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "phones retrieved successfully",
			"data":   phones,
		})
	}
}

// FindPhoneByID godoc
// @Summary Busca telefone por ID
// @Description Busca um telefone pelo ID informado como query parameter.
// @Tags phones
// @Produce json
// @Param id query uint64 true "ID do telefone"
// @Success 200 {object} map[string]interface{} "Telefone encontrado com sucesso"
// @Failure 400 {object} map[string]string "ID inválido ou não informado"
// @Failure 404 {object} map[string]string "Telefone não encontrado"
// @Router /phones/id [get]
func (ctl *PhoneController) FindByIDController() gin.HandlerFunc {
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

		phone, err := ctl.phoneUseCase.FindByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "phone found successfully",
			"data":   phone,
		})
	}
}

// FindPhoneByEmployeeID godoc
// @Summary Busca telefones por funcionário
// @Description Retorna os telefones vinculados ao funcionário informado pelo ID.
// @Tags phones
// @Produce json
// @Param employee_id query uint64 true "ID do funcionário"
// @Success 200 {object} map[string]interface{} "Telefones encontrados com sucesso"
// @Failure 400 {object} map[string]string "ID do funcionário inválido ou não informado"
// @Failure 404 {object} map[string]string "Nenhum telefone encontrado para o funcionário"
// @Router /phones/employee_id [get]
func (ctl *PhoneController) FindByEmployeeIDController() gin.HandlerFunc {
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

		phones, err := ctl.phoneUseCase.FindByEmployeeID(employeeId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "phones found successfully",
			"data":   phones,
		})
	}
}
