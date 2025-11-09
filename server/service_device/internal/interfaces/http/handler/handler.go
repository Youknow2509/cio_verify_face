package handler

import (
	"strconv"

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
	GetDeviceToken(c *gin.Context)
	RefreshDeviceToken(c *gin.Context)
	UpdateStatusDevice(c *gin.Context)
}

/**
 * Handler struct
 */
type Handler struct{}

// UpdateStatusDevice implements iHandler.
// @Summary      Update status device
// @Description  Update status device
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        request   body dto.UpdateStatusDeviceRequest  true  "Request body update device status"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/device/status [post]
func (h *Handler) UpdateStatusDevice(c *gin.Context) {
	// Get req and parse
	var req dto.UpdateStatusDeviceRequest
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	deviceUuid, err := uuidShared.ParseUUID(req.DeviceId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device_id")
		return
	}
	// Call to application handler
	errReq := applicationService.GetDeviceService().UpdateStatusDevice(
		c,
		&applicationModel.UpdateStatusDeviceInput{
			DeviceId: deviceUuid,
			Status:   req.Status,
			//
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
			CompanyId:   companyUuid,
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
	response.SuccessResponse(c, 200, "Update status device success")
}

// RefreshDeviceToken implements iHandler.
// @Summary      Refresh device access token
// @Description  Refresh device access token
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        device_id   path string  true  "Device ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/device/token/refresh/{device_id} [post]
func (h *Handler) RefreshDeviceToken(c *gin.Context) {
	// Get id device from path
	idDeviceStr := c.Param("device_id")
	idDevice, err := uuidShared.ParseUUID(idDeviceStr)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device ID")
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	// Call to application handler
	resp, errReq := applicationService.GetDeviceService().RefreshDeviceToken(
		c,
		&applicationModel.RefreshDeviceTokenInput{
			DeviceId:    idDevice,
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
			CompanyId:   companyUuid,
		},
	)
	if errReq != nil {
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, resp)
}

// GetDeviceToken implements iHandler.
// @Summary      Get device access token
// @Description  Get device access token
// @Tags         Core Device
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        device_id   path string  true  "Device ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/device/token/{device_id} [get]
func (h *Handler) GetDeviceToken(c *gin.Context) {
	// Get id device from path
	idDeviceStr := c.Param("device_id")
	idDevice, err := uuidShared.ParseUUID(idDeviceStr)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device ID")
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	// Call to application handler
	resp, errReq := applicationService.GetDeviceService().GetDeviceToken(
		c,
		&applicationModel.GetDeviceTokenInput{
			DeviceId:    idDevice,
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
			CompanyId:   companyUuid,
		},
	)
	if errReq != nil {
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, resp)
}

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
// @Router       /v1/device [post]
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
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, validationErrors.Error())
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	var companyUuidReq uuid.UUID
	if req.CompanyId != "" {
		companyUuidReq, _ = uuidShared.ParseUUID(req.CompanyId)
	}
	// Call to application handler
	service := applicationService.GetDeviceService()
	resp, errReq := service.CreateNewDevice(
		c,
		&applicationModel.CreateNewDeviceInput{
			CompanyIdReq: companyUuidReq,
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
			CompanyId:   companyUuid,
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
// @Param        device_id   path string  true  "Device ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/device/{device_id} [delete]
func (h *Handler) DeleteDeviceById(c *gin.Context) {
	// Get id device from path
	idDeviceStr := c.Param("device_id")
	idDevice, err := uuidShared.ParseUUID(idDeviceStr)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device ID")
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
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
			CompanyId:   companyUuid,
		},
	); err != nil {
		if err.ErrorSystem != nil {
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
// @Param        device_id   path string  true  "Device ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/device/{device_id} [get]
func (h *Handler) GetDeviceById(c *gin.Context) {
	// Get id device from path
	idDeviceStr := c.Param("device_id")
	idDevice, err := uuidShared.ParseUUID(idDeviceStr)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid device ID")
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
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
			CompanyId:   companyUuid,
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
// @Param        page    query     string  false  "Page number"  Format(int)
// @Param        size    query     string  false  "Page size"  Format(int)
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/device [get]
func (h *Handler) GetListDevices(c *gin.Context) {
	// Get query params
	page := c.DefaultQuery("page", strconv.Itoa(constants.PageDefault))
	size := c.DefaultQuery("size", strconv.Itoa(constants.SizeDefault))
	companyIdReq := c.DefaultQuery("company_id", "")
	// validate query params
	var err error
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid page")
		return
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil || sizeInt <= 0 || sizeInt > 100 {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid size")
		return
	}
	// Get data auth from token
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	var companyUuidReq uuid.UUID
	if companyIdReq != "" {
		companyUuidReq, err = uuidShared.ParseUUID(companyIdReq)
		if err != nil {
			response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Invalid company_id")
			return
		}
	}
	// Call to application handler
	resp, errReq := applicationService.GetDeviceService().GetListDevices(
		c,
		&applicationModel.ListDevicesInput{
			CompanyIdReq: companyUuidReq,
			Page:         pageInt,
			Size:         sizeInt,
			//
			UserId:      userUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			SessionId:   sessionUuid,
			CompanyId:   companyUuid,
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
// @Router       /v1/device [put]
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
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
			CompanyId:   companyUuid,
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
// @Router       /v1/device/location [post]
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
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
			CompanyId:   companyUuid,
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
// @Router       /v1/device/name [post]
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
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
			CompanyId:   companyUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
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
// @Router       /v1/device/info [post]
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Parse uuid
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
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
			CompanyId:   companyUuid,
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
