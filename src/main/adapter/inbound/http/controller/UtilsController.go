package controller

import (
	"net/http"

	"asset-control/src/main/usecase"

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
