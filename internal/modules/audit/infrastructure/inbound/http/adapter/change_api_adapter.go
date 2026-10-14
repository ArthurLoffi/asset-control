package adapter

import (
	"asset-control/internal/core/httpx"
	"asset-control/internal/modules/audit/application/usecase"
	"asset-control/internal/modules/audit/domain/dto"

	"github.com/gin-gonic/gin"

	// Blank imports let swag resolve the response types named in annotations.
	_ "asset-control/internal/core/dto"
	_ "asset-control/internal/core/exception"
)

type ChangeAPIAdapter struct {
	findAllChangesUseCase *usecase.FindAllChangesUseCase
}

func NewChangeAPIAdapter(findAllChangesUseCase *usecase.FindAllChangesUseCase) *ChangeAPIAdapter {
	return &ChangeAPIAdapter{findAllChangesUseCase: findAllChangesUseCase}
}

func (a *ChangeAPIAdapter) RegisterRoutes(router gin.IRouter) {
	router.GET("/changes", a.findAll)
}

// findAll godoc
// @Summary      Lista alterações do inventário
// @Description  Resolve o nome do editor e devolve só os campos escalares alterados, sem segredos nem snapshots.
// @Tags         changes
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        direction  query  string  false  "DESC para as mais recentes primeiro"  Enums(ASC, DESC)
// @Success      200  {object}  dto.PageDoc{content=[]dto.ChangeResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /changes [get]
func (a *ChangeAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.ChangeSortableFields, a.findAllChangesUseCase.Execute)
}
