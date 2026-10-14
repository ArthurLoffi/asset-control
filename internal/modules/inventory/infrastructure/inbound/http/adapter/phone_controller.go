package adapter

import (
	"asset-control/internal/core/httpx"
	"net/http"

	"asset-control/internal/modules/inventory/application/usecase"
	"asset-control/internal/modules/inventory/domain/dto"

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
// @Failure 400 {object} exception.ErrorResponse "Dados inválidos"
// @Failure 500 {object} exception.ErrorResponse "Erro interno do servidor"
// @Router /phones [post]
func (ctl *PhoneController) CreatePhoneController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var phone dto.CreatePhoneDto

		if err := httpx.BindJSON(c, &phone); err != nil {
			httpx.WriteError(c, err)
			return
		}

		created, err := ctl.phoneUseCase.CreatePhone(c.Request.Context(), phone)
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

// UpdatePhone godoc
// @Summary Atualiza um telefone
// @Description Atualiza os dados de um telefone existente através do ID informado como query parameter.
// @Tags phones
// @Accept json
// @Produce json
// @Param id query uint64 true "ID do telefone"
// @Param phone body dto.CreatePhoneDto true "Dados atualizados do telefone"
// @Success 200 {object} map[string]interface{} "Telefone atualizado com sucesso"
// @Failure 400 {object} exception.ErrorResponse "ID ou dados inválidos"
// @Failure 404 {object} exception.ErrorResponse "Telefone não encontrado"
// @Router /phones [put]
func (ctl *PhoneController) UpdatePhoneController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := httpx.QueryID(c, "id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		var phone dto.CreatePhoneDto
		if err := httpx.BindJSON(c, &phone); err != nil {
			httpx.WriteError(c, err)
			return
		}

		updated, err := ctl.phoneUseCase.UpdatePhone(c.Request.Context(), id, phone)
		if err != nil {
			httpx.WriteError(c, err)
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
// @Failure 500 {object} exception.ErrorResponse "Erro interno do servidor"
// @Router /phones [get]
func (ctl *PhoneController) GetAllPhonesController() gin.HandlerFunc {
	return func(c *gin.Context) {
		phones, err := ctl.phoneUseCase.GetAllPhones(c.Request.Context())
		if err != nil {
			httpx.WriteError(c, err)
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
// @Failure 400 {object} exception.ErrorResponse "ID inválido ou não informado"
// @Failure 404 {object} exception.ErrorResponse "Telefone não encontrado"
// @Router /phones/id [get]
func (ctl *PhoneController) FindByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := httpx.QueryID(c, "id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		phone, err := ctl.phoneUseCase.FindByID(c.Request.Context(), id)
		if err != nil {
			httpx.WriteError(c, err)
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
// @Failure 400 {object} exception.ErrorResponse "ID do funcionário inválido ou não informado"
// @Failure 404 {object} exception.ErrorResponse "Nenhum telefone encontrado para o funcionário"
// @Router /phones/employee_id [get]
func (ctl *PhoneController) FindByEmployeeIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeId, err := httpx.QueryID(c, "employee_id")
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		phones, err := ctl.phoneUseCase.FindByEmployeeID(c.Request.Context(), employeeId)
		if err != nil {
			httpx.WriteError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "phones found successfully",
			"data":   phones,
		})
	}
}
