package handler

import (
	gin "github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	// applicationModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	// applicationService "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
	// constants "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/constants"
	// dto "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/interfaces/dto"
	// response "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/interfaces/response"
	// contextShared "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/context"
	// uuidShared "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/uuid"
)

/**
 * Interface handler for http
 */
type iHandler interface {
	// For shift
	GetInfoShiftUser(*gin.Context)
	CreateShiftUser(*gin.Context)
	GetDetailShiftUser(*gin.Context)
	EditShiftUser(*gin.Context)
	DeleteShiftUser(*gin.Context)

	// For schedule
	GetInfoScheduleUser(*gin.Context)
	CreateScheduleUser(*gin.Context)
	GetDetailScheduleUser(*gin.Context)
	EditScheduleUser(*gin.Context)
	DeleteScheduleUser(*gin.Context)
}

/**
 * Handler struct
 */
type Handler struct{}

// CreateScheduleUser implements iHandler.
func (h *Handler) CreateScheduleUser(*gin.Context) {
	panic("unimplemented")
}

// CreateShiftUser implements iHandler.
func (h *Handler) CreateShiftUser(*gin.Context) {
	panic("unimplemented")
}

// DeleteScheduleUser implements iHandler.
func (h *Handler) DeleteScheduleUser(*gin.Context) {
	panic("unimplemented")
}

// DeleteShiftUser implements iHandler.
func (h *Handler) DeleteShiftUser(*gin.Context) {
	panic("unimplemented")
}

// EditScheduleUser implements iHandler.
func (h *Handler) EditScheduleUser(*gin.Context) {
	panic("unimplemented")
}

// EditShiftUser implements iHandler.
func (h *Handler) EditShiftUser(*gin.Context) {
	panic("unimplemented")
}

// GetDetailScheduleUser implements iHandler.
func (h *Handler) GetDetailScheduleUser(*gin.Context) {
	panic("unimplemented")
}

// GetDetailShiftUser implements iHandler.
func (h *Handler) GetDetailShiftUser(*gin.Context) {
	panic("unimplemented")
}

// GetInfoScheduleUser implements iHandler.
func (h *Handler) GetInfoScheduleUser(*gin.Context) {
	panic("unimplemented")
}

// GetInfoShiftUser implements iHandler.
func (h *Handler) GetInfoShiftUser(*gin.Context) {
	panic("unimplemented")
}

/**
 * New handler and impl interface
 */
func NewHandler() iHandler {
	return &Handler{}
}
