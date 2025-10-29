package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/http/handler"
)

// =======================================
// Health route
// =======================================
type HealthRoute struct {
}

// Init Health route
func (r *HealthRoute) InitHealthRoute(c *gin.RouterGroup) {
	gr := c.Group("/health")
	gr.GET("", handler.GetHealthHandler().Base)
	gr.GET("/details", handler.GetHealthHandler().SystemDetails)
}
