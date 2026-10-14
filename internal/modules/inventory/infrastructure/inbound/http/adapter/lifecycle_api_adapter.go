package adapter

import (
	"context"

	"asset-control/internal/core/dto"
	"asset-control/internal/core/httpx"

	"github.com/gin-gonic/gin"

	// Blank import lets swag resolve the error type named in annotations.
	_ "asset-control/internal/core/exception"
)

// deactivate godoc
// @Summary      Desativa um registro (soft delete)
// @Description  Registra deactivated_at uma única vez; repetir devolve o mesmo horário. Desativar um funcionário remove-o dos equipamentos ativos.
// @Tags         deactivation
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  dto.DeactivationResponseDTO
// @Failure      404  {object}  exception.ErrorResponse  "XXX_001"
// @Router       /employees/{id} [delete]
// @Router       /machines/{id} [delete]
// @Router       /phones/{id} [delete]
// @Router       /monitors/{id} [delete]
// @Router       /printers/{id} [delete]
// @Router       /emails/{id} [delete]
func deactivate(execute func(context.Context, uint64) (dto.DeactivationResponseDTO, error)) gin.HandlerFunc {
	return func(c *gin.Context) { httpx.Get(c, execute) }
}

// findInactive godoc
// @Summary      Lista o histórico de desativação
// @Description  Retorna só id e deactivated_at, ordenados por id, separado do inventário ativo.
// @Tags         deactivation
// @Produce      json
// @Param        page       query  int     false  "Página, começando em 0"  minimum(0)
// @Param        page_size  query  int     false  "Itens por página (0 = 20)"  minimum(0) maximum(100)
// @Param        direction  query  string  false  "Direção da ordenação"  Enums(ASC, DESC)
// @Success      200  {object}  dto.PageDoc{content=[]dto.DeactivationResponseDTO}
// @Router       /employees/inactive [get]
// @Router       /machines/inactive [get]
// @Router       /phones/inactive [get]
// @Router       /monitors/inactive [get]
// @Router       /printers/inactive [get]
// @Router       /emails/inactive [get]
func findInactive(execute func(context.Context, dto.PaginatedRequestDTO) (dto.PaginatedResponseDTO[dto.DeactivationResponseDTO], error)) gin.HandlerFunc {
	return func(c *gin.Context) { httpx.List(c, []string{"id"}, execute) }
}
