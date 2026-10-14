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

type PhoneUseCases = resource.UseCases[dto.PhoneRequestDTO, dto.PhoneInputDTO, dto.PhoneResponseDTO]

type PhoneAPIAdapter struct{ useCases PhoneUseCases }

func NewPhoneAPIAdapter(useCases PhoneUseCases) *PhoneAPIAdapter {
	return &PhoneAPIAdapter{useCases: useCases}
}

func (a *PhoneAPIAdapter) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/phones")
	group.GET("", a.findAll)
	group.GET("/inactive", findInactive(a.useCases.FindInactive.Execute))
	group.GET("/:id", a.findByID)
	group.POST("", a.create)
	group.PUT("/:id", a.update)
	group.DELETE("/:id", deactivate(a.useCases.Deactivate.Execute))
}

// findAll godoc
// @Summary      Lista celulares ativos
// @Tags         phones
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Param        sort_by    query  string  false  "Campo de ordenação; o id desempata"
// @Param        employee_id query  []int   false  "IDs de funcionários (repetir o parâmetro; qualquer um)"
// @Param        unassigned query  bool    false  "Somente sem funcionário"
// @Success      200  {object}  dto.PageDoc{content=[]dto.PhoneResponseDTO}
// @Failure      400  {object}  exception.ErrorResponse  "CORE_001"
// @Router       /phones [get]
func (a *PhoneAPIAdapter) findAll(c *gin.Context) {
	httpx.List(c, dto.PhoneSortableFields, a.useCases.FindAll.Execute)
}

// findByID godoc
// @Summary      Busca celular ativo por ID
// @Tags         phones
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.PhoneResponseDTO
// @Failure      400  {object}  exception.ErrorResponse  "CORE_002"
// @Failure      404  {object}  exception.ErrorResponse  "PHONE_001"
// @Router       /phones/{id} [get]
func (a *PhoneAPIAdapter) findByID(c *gin.Context) {
	httpx.Get(c, a.useCases.FindByID.Execute)
}

// create godoc
// @Summary      Cria celular
// @Description  ID e created_at são gerados pelo servidor.
// @Tags         phones
// @Accept       json
// @Produce      json
// @Param        body  body      dto.PhoneInputDTO  true  "Dados"
// @Success      201   {object}  dto.PhoneResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      409   {object}  exception.ErrorResponse  "CORE_004"
// @Router       /phones [post]
func (a *PhoneAPIAdapter) create(c *gin.Context) {
	httpx.Create(c, a.useCases.Create.Execute)
}

// update godoc
// @Summary      Atualiza celular ativo
// @Description  Substitui os campos editáveis; preserva id, created_at e campos desconhecidos. Registra auditoria.
// @Tags         phones
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "ID"
// @Param        body  body      dto.PhoneInputDTO  true  "Dados"
// @Success      200   {object}  dto.PhoneResponseDTO
// @Failure      400   {object}  exception.ErrorResponse  "CORE_001 ou CORE_002"
// @Failure      404   {object}  exception.ErrorResponse  "PHONE_001"
// @Router       /phones/{id} [put]
func (a *PhoneAPIAdapter) update(c *gin.Context) {
	httpx.Update(c, a.useCases.Update.Execute)
}
