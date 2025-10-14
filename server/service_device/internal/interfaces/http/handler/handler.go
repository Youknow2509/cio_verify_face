package handler

import "github.com/gin-gonic/gin"

/**
 * Interface handler for http
 */
type iHandler interface {
	GetListDevices(c *gin.Context)
	CreateNewDevice(c *gin.Context)
	GetDeviceById(c *gin.Context)
	UpdateDeviceById(c *gin.Context)
	DeleteDeviceById(c *gin.Context)
}

/**
 * Handler struct
 */
type Handler struct{}

// CreateNewDevice implements iHandler.
func (h *Handler) CreateNewDevice(c *gin.Context) {
	panic("unimplemented")
}

// DeleteDeviceById implements iHandler.
func (h *Handler) DeleteDeviceById(c *gin.Context) {
	panic("unimplemented")
}

// GetDeviceById implements iHandler.
func (h *Handler) GetDeviceById(c *gin.Context) {
	panic("unimplemented")
}

// GetListDevices implements iHandler.
func (h *Handler) GetListDevices(c *gin.Context) {
	panic("unimplemented")
}

// UpdateDeviceById implements iHandler.
func (h *Handler) UpdateDeviceById(c *gin.Context) {
	panic("unimplemented")
}

/**
 * New handler and impl interface
 */
func NewHandler() iHandler {
	return &Handler{}
}
