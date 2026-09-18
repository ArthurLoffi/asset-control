package controller

import (
	"asset-control/src/main/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DeactivationController struct{ usecase *usecase.DeactivationUseCase }

func NewDeactivationController(uc *usecase.DeactivationUseCase) *DeactivationController {
	return &DeactivationController{usecase: uc}
}

// Deactivate godoc
// @Summary Deactivate an inventory record
// @Tags deactivation
// @Produce json
// @Param id query uint64 true "Record ID"
// @Success 200 {object} dto.DeactivationResponseDto
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /employees/deactivate [post]
// @Router /machines/deactivate [post]
// @Router /phones/deactivate [post]
// @Router /monitors/deactivate [post]
// @Router /printers/deactivate [post]
// @Router /emails/deactivate [post]
func (ctl *DeactivationController) Deactivate(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseID(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}
		result, err := ctl.usecase.Deactivate(c.Request.Context(), resource, id)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deactivated", "data": result})
	}
}

// ListInactive godoc
// @Summary List deactivation history separately from active inventory
// @Tags deactivation
// @Produce json
// @Param after query uint64 false "Last ID from previous page (default 0)"
// @Param limit query int false "Page size (1-100, default 50)"
// @Success 200 {object} dto.InactivePageResponseDto
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /employees/inactive [get]
// @Router /machines/inactive [get]
// @Router /phones/inactive [get]
// @Router /monitors/inactive [get]
// @Router /printers/inactive [get]
// @Router /emails/inactive [get]
func (ctl *DeactivationController) ListInactive(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		after, err := strconv.ParseUint(c.DefaultQuery("after", "0"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'after' value"})
			return
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'limit' value"})
			return
		}
		page, err := ctl.usecase.ListInactive(c.Request.Context(), resource, after, limit)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "inactive records retrieved", "data": page})
	}
}
