package adapter

import (
	"asset-control/internal/core/httpx"
	"asset-control/internal/core/resource"
	"asset-control/internal/modules/inventory/domain/dto"

	"github.com/gin-gonic/gin"

	// Blank imports let swag resolve the response types named in annotations.
	_ "asset-control/internal/core/dto"
	_ "asset-control/internal/core/exception"
)

type EmployeeUseCases = resource.UseCases[dto.EmployeeRequestDTO, dto.EmployeeInputDTO, dto.EmployeeResponseDTO]

type EmployeeAPIAdapter struct{ useCases EmployeeUseCases }

func NewEmployeeAPIAdapter(useCases EmployeeUseCases) *EmployeeAPIAdapter {
	return &EmployeeAPIAdapter{useCases: useCases}
}

func (a *EmployeeAPIAdapter) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/employees")
	group.GET("", a.findAll)
	group.GET("/inactive", findInactive(a.useCases.FindInactive.Execute))
	group.GET("/:id", a.findByID)
	group.POST("", a.create)
	group.PUT("/:id", a.update)
	group.DELETE("/:id", deactivate(a.useCases.Deactivate.Execute))
}

// findAll godoc
// @Summary      Lista funcionários ativos
// @Tags         employees
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Param        sort_by    query  string  false  "Campo de ordenação; o id desempata"
// @Param        name       query  string  false  "Nome (parcial, sem diferenciar maiúsculas)"
// @Param        branch_id  query  int     false  "ID da filial"
// @Param        department_id query  int     false  "ID do departamento"
// @Success      200  {object}  dto.PageDoc{content=[]dto.EmployeeResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /employees [get]
func (a *EmployeeAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.EmployeeSortableFields, a.useCases.FindAll.Execute)
}

// findByID godoc
// @Summary      Busca funcionário ativo por ID
// @Tags         employees
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.EmployeeResponseDTO
// @Failure      400  {object}  exception.ErrorResponse  "CORE_002"
// @Failure      404  {object}  exception.ErrorResponse  "EMPLOYEE_001"
// @Router       /employees/{id} [get]
func (a *EmployeeAPIAdapter) findByID(c *gin.Context) {
	httpx.Get(c, a.useCases.FindByID.Execute)
}

// create godoc
// @Summary      Cria funcionário
// @Description  ID e created_at são gerados pelo servidor.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        body  body      dto.EmployeeInputDTO  true  "Dados"
// @Success      201   {object}  dto.EmployeeResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      409   {object}  exception.ErrorResponse  "CORE_004"
// @Router       /employees [post]
func (a *EmployeeAPIAdapter) create(c *gin.Context) {
	httpx.Create(c, a.useCases.Create.Execute)
}

// update godoc
// @Summary      Atualiza funcionário ativo
// @Description  Substitui os campos editáveis; preserva id, created_at e campos desconhecidos. Registra auditoria.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "ID"
// @Param        body  body      dto.EmployeeInputDTO  true  "Dados"
// @Success      200   {object}  dto.EmployeeResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      404   {object}  exception.ErrorResponse  "EMPLOYEE_001"
// @Router       /employees/{id} [put]
func (a *EmployeeAPIAdapter) update(c *gin.Context) {
	httpx.Update(c, a.useCases.Update.Execute)
}
