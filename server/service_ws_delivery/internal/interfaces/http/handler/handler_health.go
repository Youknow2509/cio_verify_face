package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
)

// ===========================================
// Health check handlers
// ===========================================
type HealthHandler struct{}

// Get Health check handlers
func GetHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Base(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "running"})
}

func (h *HealthHandler) SystemDetails(c *gin.Context) {
	data := applicationService.GetHealthCheckService().SystemDetails(
		c,
		&model.SystemDetailsInput{
			ClientIp: c.ClientIP(),
		},
	)
	if data == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get system details"})
		return
	}
	c.JSON(http.StatusOK, data)
}
