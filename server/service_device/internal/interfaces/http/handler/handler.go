package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/interfaces/dto"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/interfaces/response"
	contextShared "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/context"
	uuidShared "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/uuid"
)

/**
 * Interface handler for http
 */
type iHandler interface {
	GetListDevices(c *gin.Context)
	CreateNewDevice(c *gin.Context)
	GetDeviceById(c *gin.Context)
	UpdateDeviceById(c *gin.Context)
	DeleteDeviceById(c *gin.Context)
	UpdateLocationDevice(c *gin.Context)
	UpdateNameDevice(c *gin.Context)
	UpdateInfoDevice(c *gin.Context)
}

/**
 * Handler struct
 */
type Handler struct{}

// CreateNewDevice implements iHandler.
// @Summary      Create new device
// @Description  Create new device
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        request   body dto.CreateDeviceRequest  true  "Request body create device"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices [post]
func (h *Handler) CreateNewDevice(c *gin.Context) {
	// Get req and parse
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrorCodeBindRequest, "Invalid request body")
		return
	}
	// Validate req
	validateMiddleware, ok := c.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	err := validate.Struct(req)
	validationErrors := err.(validator.ValidationErrors)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, validationErrors.Error())
		return
	}
	companyId, err := uuidShared.ParseUUID(req.CompanyId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid company_id")
		return
	}
	var locationId uuid.UUID
	if req.LocationId != "" {
		locationId, err = uuidShared.ParseUUID(req.LocationId)
		if err != nil {
			response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid location_id")
			return
		}
	}
	// Get data auth from token
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call to application handler
	service := applicationService.GetDeviceService()
	resp, errReq := service.CreateNewDevice(
		c,
		&applicationModel.CreateNewDeviceInput{
			CompanyId:    companyId,
			LocationId:   locationId,
			DeviceName:   req.DeviceName,
			Address:      req.Address,
			DeviceType:   req.DeviceType,
			SerialNumber: req.SerialNumber,
			MacAddress:   req.MacAddress,
			//
			UserId:      userUuid,
			Role:        userRole,
			SessionId:   sessionUuid,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
		},
	)
	if errReq != nil {
		if errReq.ErrorClient == "" {
			response.ErrorResponse(c, 500, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, resp)
}

// DeleteDeviceById implements iHandler.
// @Summary      Delete device by ID
// @Description  Delete device by ID
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        id   path string  true  "Device ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices/{id} [delete]
func (h *Handler) DeleteDeviceById(c *gin.Context) {
	// Get id device from path
	idDeviceStr := c.Param("id")
	idDevice, err := uuidShared.ParseUUID(idDeviceStr)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device ID")
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call to application handler
	if err := applicationService.GetDeviceService().DeleteDeviceById(
		c,
		&applicationModel.DeleteDeviceInput{
			DeviceId:    idDevice,
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
		},
	); err != nil {
		if err.ErrorClient == "" {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, response.ErrorCodeDeleteDevice, err.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Delete device success")
}

// GetDeviceById implements iHandler.
// @Summary      Delete device by ID
// @Description  Delete device by ID
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        id   path string  true  "Device ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices/{id} [delete]
func (h *Handler) GetDeviceById(c *gin.Context) {
	// Get id device from path
	idDeviceStr := c.Param("id")
	idDevice, err := uuidShared.ParseUUID(idDeviceStr)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device ID")
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call to application handler
	resp, errReq := applicationService.GetDeviceService().GetDeviceById(
		c,
		&applicationModel.GetDeviceByIdInput{
			DeviceId:    idDevice,
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
		},
	)
	if errReq != nil {
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, resp)
}

// GetListDevices implements iHandler.
// @Summary      Get list of devices
// @Description  Get list of devices
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        request   body dto.ListDevicesRequest  true  "Request body list devices"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices [get]
func (h *Handler) GetListDevices(c *gin.Context) {
	// Get req and parse
	var req dto.ListDevicesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorResponse(c, response.ErrorCodeBindRequest, "Invalid request body")
		return
	}
	// Validate req
	validateMiddleware, ok := c.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	err := validate.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, validationErrors.Error())
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	companyUuid, err := uuidShared.ParseUUID(req.CompanyId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid company_id")
		return
	}
	// Call to application handler
	resp, errReq := applicationService.GetDeviceService().GetListDevices(
		c,
		&applicationModel.ListDevicesInput{
			CompanyId: companyUuid,
			Page:      req.Page,
			Size:      req.Size,
			//
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorClient == "" {
			response.ErrorResponse(c, 500, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, resp)
}

// UpdateDeviceById implements iHandler.
// @Summary      Update device by ID
// @Description  Update device by ID
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        request   body dto.UpdateDeviceRequest  true  "Request body update device"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices [put]
func (h *Handler) UpdateDeviceById(c *gin.Context) {
	// Get req and parse
	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrorCodeBindRequest, "Invalid request body")
		return
	}
	// Validate req
	validateMiddleware, ok := c.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	err := validate.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, validationErrors.Error())
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	locationUuid, err := uuidShared.ParseUUID(req.LocationId)
	if err != nil && req.LocationId != "" {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid location_id")
		return
	}
	// Call to application handler
	resp, errReq := applicationService.GetDeviceService().UpdateDeviceById(
		c,
		&applicationModel.UpdateDeviceInput{
			LocationId:   locationUuid,
			DeviceName:   req.DeviceName,
			Address:      req.Address,
			DeviceType:   req.DeviceType,
			SerialNumber: req.SerialNumber,
			MacAddress:   req.MacAddress,
			Status:       req.Status,
			//
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorClient == "" {
			response.ErrorResponse(c, 500, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, resp)
}

// Update location
// @Summary      Update location
// @Description  Update location
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        request   body dto.UpdateLocationDeviceRequest  true  "Request body update device"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices/location [post]
func (h *Handler) UpdateLocationDevice(c *gin.Context) {
	// Get req and parse
	var req dto.UpdateLocationDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrorCodeBindRequest, "Invalid request body")
		return
	}
	// Validate req
	validateMiddleware, ok := c.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	err := validate.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, validationErrors.Error())
		return
	}
	// Get ession
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	deviceUuid, err := uuidShared.ParseUUID(req.DeviceId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device_id")
		return
	}
	newLocationUuid, err := uuidShared.ParseUUID(req.NewLocationId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid location_id")
		return
	}
	// Call to application handler
	errReq := applicationService.GetDeviceService().UpdateLocationDevice(
		c,
		&applicationModel.UpdateLocationDeviceInput{
			DeviceId:      deviceUuid,
			NewLocationId: newLocationUuid,
			NewAddress:    req.NewAddress,
			//
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorClient == "" {
			response.ErrorResponse(c, 500, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Update location device success")
}

// Update name device
// @Description  Update name device
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        request   body dto.UpdateNameDeviceRequest  true  "Request body update device"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices/name [post]
func (h *Handler) UpdateNameDevice(c *gin.Context) {
	// Get req and parse
	var req dto.UpdateNameDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrorCodeBindRequest, "Invalid request body")
		return
	}
	// Validate req
	validateMiddleware, ok := c.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	err := validate.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, validationErrors.Error())
		return
	}
	// Get ession
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	deviceUuid, err := uuidShared.ParseUUID(req.DeviceId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device_id")
		return
	}
	// Call to application handler
	errReq := applicationService.GetDeviceService().UpdateNameDevice(
		c,
		&applicationModel.UpdateNameDeviceInput{
			DeviceId:    deviceUuid,
			NewName:     req.NewName,
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
		},
	)
	if errReq != nil {
		if errReq.ErrorClient == "" {
			response.ErrorResponse(c, 500, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Update name device success")
}

// UpdateInfoDevice implements iHandler.
// @Description  Update info device
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        request   body dto.UpdateInfoDeviceRequest  true  "Request body update device"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/devices/info [post]
func (h *Handler) UpdateInfoDevice(c *gin.Context) {
	// Get req and parse
	var req dto.UpdateInfoDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrorCodeBindRequest, "Invalid request body")
		return
	}
	// Validate req
	validateMiddleware, ok := c.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	err := validate.Struct(req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, validationErrors.Error())
		return
	}
	// Get ession
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	deviceUuid, err := uuidShared.ParseUUID(req.DeviceId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device_id")
		return
	}
	// Call to application handler
	errReq := applicationService.GetDeviceService().UpdateInfoDevice(
		c,
		&applicationModel.UpdateInfoDeviceInput{
			DeviceId:        deviceUuid,
			NewDeviceType:   req.NewDeviceType,
			NewSerialNumber: req.NewSerialNumber,
			NewMacAddress:   req.NewMacAddress,
			//
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
		},
	)
	if errReq != nil {
		if errReq.ErrorClient == "" {
			response.ErrorResponse(c, 500, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Update info device success")
}

/**
 * New handler and impl interface
 */
func NewHandler() iHandler {
	return &Handler{}
}
