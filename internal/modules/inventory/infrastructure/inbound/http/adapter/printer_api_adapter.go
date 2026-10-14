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

type PrinterUseCases = resource.UseCases[dto.PrinterRequestDTO, dto.PrinterInputDTO, dto.PrinterResponseDTO]

type PrinterAPIAdapter struct{ useCases PrinterUseCases }

func NewPrinterAPIAdapter(useCases PrinterUseCases) *PrinterAPIAdapter {
	return &PrinterAPIAdapter{useCases: useCases}
}

func (a *PrinterAPIAdapter) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/printers")
	group.GET("", a.findAll)
	group.GET("/inactive", findInactive(a.useCases.FindInactive.Execute))
	group.GET("/:id", a.findByID)
	group.POST("", a.create)
	group.PUT("/:id", a.update)
	group.DELETE("/:id", deactivate(a.useCases.Deactivate.Execute))
}

// findAll godoc
// @Summary      Lista impressoras ativos
// @Tags         printers
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Param        sort_by    query  string  false  "Campo de ordenação; o id desempata"
// @Param        branch_id  query  int     false  "ID da filial"
// @Param        departments_id query  int     false  "ID do departamento"
// @Success      200  {object}  dto.PageDoc{content=[]dto.PrinterResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /printers [get]
func (a *PrinterAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.PrinterSortableFields, a.useCases.FindAll.Execute)
}

// findByID godoc
// @Summary      Busca impressora ativo por ID
// @Tags         printers
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.PrinterResponseDTO
// @Failure      400  {object}  exception.ErrorResponse  "CORE_002"
// @Failure      404  {object}  exception.ErrorResponse  "PRINTER_001"
// @Router       /printers/{id} [get]
func (a *PrinterAPIAdapter) findByID(c *gin.Context) {
	httpx.Get(c, a.useCases.FindByID.Execute)
}

// create godoc
// @Summary      Cria impressora
// @Description  ID e created_at são gerados pelo servidor.
// @Tags         printers
// @Accept       json
// @Produce      json
// @Param        body  body      dto.PrinterInputDTO  true  "Dados"
// @Success      201   {object}  dto.PrinterResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      409   {object}  exception.ErrorResponse  "CORE_004"
// @Router       /printers [post]
func (a *PrinterAPIAdapter) create(c *gin.Context) {
	httpx.Create(c, a.useCases.Create.Execute)
}

// update godoc
// @Summary      Atualiza impressora ativo
// @Description  Substitui os campos editáveis; preserva id, created_at e campos desconhecidos. Registra auditoria.
// @Tags         printers
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "ID"
// @Param        body  body      dto.PrinterInputDTO  true  "Dados"
// @Success      200   {object}  dto.PrinterResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      404   {object}  exception.ErrorResponse  "PRINTER_001"
// @Router       /printers/{id} [put]
func (a *PrinterAPIAdapter) update(c *gin.Context) {
	httpx.Update(c, a.useCases.Update.Execute)
}
