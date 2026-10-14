package adapter

import (
	"asset-control/internal/core/httpx"
	"asset-control/internal/core/resource"
	"asset-control/internal/modules/organization/domain/dto"

	"github.com/gin-gonic/gin"

	// Blank imports let swag resolve the response types named in annotations.
	_ "asset-control/internal/core/dto"
	_ "asset-control/internal/core/exception"
)

type BranchUseCases = resource.ReadUseCases[dto.BranchRequestDTO, dto.BranchResponseDTO]

type BranchAPIAdapter struct{ useCases BranchUseCases }

func NewBranchAPIAdapter(useCases BranchUseCases) *BranchAPIAdapter {
	return &BranchAPIAdapter{useCases: useCases}
}

func (a *BranchAPIAdapter) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/branches")
	group.GET("", a.findAll)
	group.GET("/:id", a.findByID)
}

// findAll godoc
// @Summary      Lista filiais
// @Tags         organization
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        sort_by    query  string  false  "Campo de ordenação"  Enums(id, name)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Param        name       query  string  false  "Nome (parcial, sem diferenciar maiúsculas)"
// @Success      200  {object}  dto.PageDoc{content=[]dto.BranchResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /branches [get]
func (a *BranchAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.OrganizationSortableFields, a.useCases.FindAll.Execute)
}

// findByID godoc
// @Summary      Busca filial por ID
// @Tags         organization
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.BranchResponseDTO
// @Failure      404  {object}  exception.ErrorResponse  "BRANCH_001"
// @Router       /branches/{id} [get]
func (a *BranchAPIAdapter) findByID(c *gin.Context) {
	httpx.Get(c, a.useCases.FindByID.Execute)
}
