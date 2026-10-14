package adapter

import (
	"asset-control/internal/core/httpx"
	"net/http"

	"asset-control/internal/modules/audit/application/usecase"
	"github.com/gin-gonic/gin"
)

type RecentChangeController struct{ usecase *usecase.RecentChangeUseCase }

func NewRecentChangeController(uc *usecase.RecentChangeUseCase) *RecentChangeController {
	return &RecentChangeController{usecase: uc}
}

// ListRecent godoc
// @Summary List the ten most recent inventory updates
// @Tags changes
// @Produce json
// @Success 200 {object} dto.RecentChangesResponseDto
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /changes/recent [get]
func (ctl *RecentChangeController) ListRecent(c *gin.Context) {
	changes, err := ctl.usecase.ListRecent(c.Request.Context())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "recent changes retrieved", "data": changes})
}
