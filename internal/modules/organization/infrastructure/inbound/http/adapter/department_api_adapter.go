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

type DepartmentUseCases = resource.ReadUseCases[dto.DepartmentRequestDTO, dto.DepartmentResponseDTO]

type DepartmentAPIAdapter struct{ useCases DepartmentUseCases }

func NewDepartmentAPIAdapter(useCases DepartmentUseCases) *DepartmentAPIAdapter {
	return &DepartmentAPIAdapter{useCases: useCases}
}

func (a *DepartmentAPIAdapter) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/departments")
	group.GET("", a.findAll)
	group.GET("/:id", a.findByID)
}

// findAll godoc
// @Summary      Lista departamentos
// @Tags         organization
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        sort_by    query  string  false  "Campo de ordenação"  Enums(id, name)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Param        name       query  string  false  "Nome (parcial, sem diferenciar maiúsculas)"
// @Success      200  {object}  dto.PageDoc{content=[]dto.DepartmentResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /departments [get]
func (a *DepartmentAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.OrganizationSortableFields, a.useCases.FindAll.Execute)
}

// findByID godoc
// @Summary      Busca departamento por ID
// @Tags         organization
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.DepartmentResponseDTO
// @Failure      404  {object}  exception.ErrorResponse  "DEPARTMENT_001"
// @Router       /departments/{id} [get]
func (a *DepartmentAPIAdapter) findByID(c *gin.Context) {
	httpx.Get(c, a.useCases.FindByID.Execute)
}
