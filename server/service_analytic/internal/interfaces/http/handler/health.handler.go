package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/dto"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	service applicationService.IAnalyticService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		service: applicationService.GetAnalyticService(),
	}
}

// HealthCheck handles GET /health
// @Summary Health check
// @Description Check service health status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	health, appErr := h.service.GetHealthCheck(c.Request.Context())
	if appErr != nil {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(health))
}
