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

type MachineUseCases = resource.UseCases[dto.MachineRequestDTO, dto.MachineInputDTO, dto.MachineResponseDTO]

type MachineAPIAdapter struct{ useCases MachineUseCases }

func NewMachineAPIAdapter(useCases MachineUseCases) *MachineAPIAdapter {
	return &MachineAPIAdapter{useCases: useCases}
}

func (a *MachineAPIAdapter) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/machines")
	group.GET("", a.findAll)
	group.GET("/inactive", findInactive(a.useCases.FindInactive.Execute))
	group.GET("/:id", a.findByID)
	group.POST("", a.create)
	group.PUT("/:id", a.update)
	group.DELETE("/:id", deactivate(a.useCases.Deactivate.Execute))
}

// findAll godoc
// @Summary      Lista máquinas ativos
// @Tags         machines
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Param        sort_by    query  string  false  "Campo de ordenação; o id desempata"
// @Param        employee_id query  []int   false  "IDs de funcionários (repetir o parâmetro; qualquer um)"
// @Param        unassigned query  bool    false  "Somente sem funcionário"
// @Success      200  {object}  dto.PageDoc{content=[]dto.MachineResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /machines [get]
func (a *MachineAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.MachineSortableFields, a.useCases.FindAll.Execute)
}

// findByID godoc
// @Summary      Busca máquina ativo por ID
// @Tags         machines
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.MachineResponseDTO
// @Failure      400  {object}  exception.ErrorResponse  "CORE_002"
// @Failure      404  {object}  exception.ErrorResponse  "MACHINE_001"
// @Router       /machines/{id} [get]
func (a *MachineAPIAdapter) findByID(c *gin.Context) {
	httpx.Get(c, a.useCases.FindByID.Execute)
}

// create godoc
// @Summary      Cria máquina
// @Description  ID e created_at são gerados pelo servidor.
// @Tags         machines
// @Accept       json
// @Produce      json
// @Param        body  body      dto.MachineInputDTO  true  "Dados"
// @Success      201   {object}  dto.MachineResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      409   {object}  exception.ErrorResponse  "CORE_004"
// @Router       /machines [post]
func (a *MachineAPIAdapter) create(c *gin.Context) {
	httpx.Create(c, a.useCases.Create.Execute)
}

// update godoc
// @Summary      Atualiza máquina ativo
// @Description  Substitui os campos editáveis; preserva id, created_at e campos desconhecidos. Registra auditoria.
// @Tags         machines
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "ID"
// @Param        body  body      dto.MachineInputDTO  true  "Dados"
// @Success      200   {object}  dto.MachineResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      404   {object}  exception.ErrorResponse  "MACHINE_001"
// @Router       /machines/{id} [put]
func (a *MachineAPIAdapter) update(c *gin.Context) {
	httpx.Update(c, a.useCases.Update.Execute)
}
