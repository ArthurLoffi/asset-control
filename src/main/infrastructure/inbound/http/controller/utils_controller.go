package controller

import (
	"net/http"

	"asset-control/src/main/application/usecase"

	"github.com/gin-gonic/gin"
)

type UtilsController struct {
	utilsUseCase *usecase.UtilsUseCase
}

func NewUtilsController(utilsUseCase *usecase.UtilsUseCase) *UtilsController {
	return &UtilsController{
		utilsUseCase: utilsUseCase,
	}
}

// GetAllBranchs godoc
// @Summary Lista todas as filiais criadas
// @Description Retorna todas as filiais criadas no banco de dados.
// @Tags branches
// @Produce json
// @Success 200 {object} map[string]interface{} "Filias encontrados com sucesso"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /utils/branches [get]
func (ctl *UtilsController) GetAllBranchs() gin.HandlerFunc {
	return func(c *gin.Context) {
		branches, err := ctl.utilsUseCase.GetAllBranchs(c.Request.Context())
		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "branches retrieved successfully",
			"data":   branches,
		})
	}
}

// GetBranchName godoc
// @Summary Mostra o nome da branch
// @Description Retorna o nome dela de acordo com o ID passado.
// @Tags branches
// @Produce json
// @Param id query uint64 true "ID da branch"
// @Success 200 {object} map[string]interface{} "Branch encontrada com sucesso"
// @Failure 404 {object} map[string]string "Branch não encontrada"
// @Router /utils/branch [get]
func (uc *UtilsController) BranchIdToDescription() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id is required",
			})
			return
		}

		id, err := parseID(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}

		description, err := uc.utilsUseCase.BranchIdToDescription(c.Request.Context(), id)
		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "branch found",
			"data": gin.H{
				"name": description,
			},
		})
	}
}

// GetAllDepartments godoc
// @Summary Lista todos os setores criados
// @Description Retorna todos os setores criados no banco de dados.
// @Tags departments
// @Produce json
// @Success 200 {object} map[string]interface{} "Setores encontrados com sucesso"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /utils/departments [get]
func (ctl *UtilsController) GetAllDepartments() gin.HandlerFunc {
	return func(c *gin.Context) {
		departments, err := ctl.utilsUseCase.GetAllDepartments(c.Request.Context())
		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "departments retrieved successfully",
			"data":   departments,
		})
	}
}

// GetDepartmentName godoc
// @Summary Mostra o nome do departamento
// @Description Retorna o nome dele de acordo com o ID passado.
// @Tags departments
// @Produce json
// @Param id query uint64 true "ID do department"
// @Success 200 {object} map[string]interface{} "Department encontrada com sucesso"
// @Failure 404 {object} map[string]string "Department não encontrada"
// @Router /utils/department [get]
func (uc *UtilsController) DepartmentIdToDescription() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id is required",
			})
			return
		}

		id, err := parseID(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}

		description, err := uc.utilsUseCase.DepartmentIdToDescription(c.Request.Context(), id)
		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "department found",
			"data": gin.H{
				"name": description,
			},
		})
	}
}
