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

type MonitorUseCases = resource.UseCases[dto.MonitorRequestDTO, dto.MonitorInputDTO, dto.MonitorResponseDTO]

type MonitorAPIAdapter struct{ useCases MonitorUseCases }

func NewMonitorAPIAdapter(useCases MonitorUseCases) *MonitorAPIAdapter {
	return &MonitorAPIAdapter{useCases: useCases}
}

func (a *MonitorAPIAdapter) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/monitors")
	group.GET("", a.findAll)
	group.GET("/inactive", findInactive(a.useCases.FindInactive.Execute))
	group.GET("/:id", a.findByID)
	group.POST("", a.create)
	group.PUT("/:id", a.update)
	group.DELETE("/:id", deactivate(a.useCases.Deactivate.Execute))
}

// findAll godoc
// @Summary      Lista monitores ativos
// @Tags         monitors
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Param        sort_by    query  string  false  "Campo de ordenação; o id desempata"
// @Param        employee_id query  []int   false  "IDs de funcionários (repetir o parâmetro; qualquer um)"
// @Param        unassigned query  bool    false  "Somente sem funcionário"
// @Success      200  {object}  dto.PageDoc{content=[]dto.MonitorResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /monitors [get]
func (a *MonitorAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.MonitorSortableFields, a.useCases.FindAll.Execute)
}

// findByID godoc
// @Summary      Busca monitor ativo por ID
// @Tags         monitors
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.MonitorResponseDTO
// @Failure      400  {object}  exception.ErrorResponse  "CORE_002"
// @Failure      404  {object}  exception.ErrorResponse  "MONITOR_001"
// @Router       /monitors/{id} [get]
func (a *MonitorAPIAdapter) findByID(c *gin.Context) {
	httpx.Get(c, a.useCases.FindByID.Execute)
}

// create godoc
// @Summary      Cria monitor
// @Description  ID e created_at são gerados pelo servidor.
// @Tags         monitors
// @Accept       json
// @Produce      json
// @Param        body  body      dto.MonitorInputDTO  true  "Dados"
// @Success      201   {object}  dto.MonitorResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      409   {object}  exception.ErrorResponse  "CORE_004"
// @Router       /monitors [post]
func (a *MonitorAPIAdapter) create(c *gin.Context) {
	httpx.Create(c, a.useCases.Create.Execute)
}

// update godoc
// @Summary      Atualiza monitor ativo
// @Description  Substitui os campos editáveis; preserva id, created_at e campos desconhecidos. Registra auditoria.
// @Tags         monitors
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "ID"
// @Param        body  body      dto.MonitorInputDTO  true  "Dados"
// @Success      200   {object}  dto.MonitorResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      404   {object}  exception.ErrorResponse  "MONITOR_001"
// @Router       /monitors/{id} [put]
func (a *MonitorAPIAdapter) update(c *gin.Context) {
	httpx.Update(c, a.useCases.Update.Execute)
}
