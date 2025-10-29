package handler

import (
	"time"

	gin "github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/constants"
	dto "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/interfaces/dto"
	response "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/interfaces/response"
	contextShared "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/context"
	uuidShared "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/uuid"
)

/**
 * Interface handler for http
 */
type iHandler interface {
	// For shift
	CreateShift(*gin.Context)
	GetDetailShift(*gin.Context)
	EditShift(*gin.Context)
	DeleteShift(*gin.Context)
	// For shift employee
	GetShiftForUserWithEffectiveDate(*gin.Context)
	EditShiftForUserWithEffectiveDate(*gin.Context)
	EnableShiftForUser(*gin.Context)
	DisableShiftForUser(*gin.Context)
	DeleteShiftForUser(*gin.Context)
	AddShiftEmployee(*gin.Context)
}

/**
 * Handler struct
 */
type Handler struct{}

// AddShiftEmployee implements iHandler.
// @Summary      Create shift information
// @Description  Create shift information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.AddShiftEmployeeReq true "Create Shift Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/add [post]
func (h *Handler) AddShiftEmployee(c *gin.Context) {
	// Get req
	var req dto.AddShiftEmployeeReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, 400, "Data input error")
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
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	employeeUuid, err := uuidShared.ParseUUID(req.EmployeeId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid employee ID")
		return
	}
	shiftUuid, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift ID")
		return
	}
	// Call service create shift employee
	reps, errReq := applicationService.GetShiftEmployeeService().AddShiftEmployee(
		c,
		&applicationModel.AddShiftEmployeeInput{
			// User info
			UserId:    userUuid,
			SessionId: sessionUuid,
			Role:      userRole,
			//
			EmployeeId:    employeeUuid,
			ShiftId:       shiftUuid,
			EffectiveFrom: time.Unix(req.EffectiveFrom, 0),
			EffectiveTo:   time.Unix(req.EffectiveTo, 0),
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, reps)
}

// CreateShift implements iHandler.
// @Summary      Create shift information
// @Description  Create shift information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.CreateShiftReq true "Create Shift Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/shift [post]
func (h *Handler) CreateShift(c *gin.Context) {
	// Binding data req
	var req dto.CreateShiftReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, 400, "Data input error")
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
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	companyId, err := uuidShared.ParseUUID(req.CompanyId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid company ID")
		return
	}
	// Call service create shift
	reps, errReq := applicationService.GetShiftService().CreateShift(
		c,
		&applicationModel.CreateShiftInput{
			// User info
			UserId:      userUuid,
			UserSession: sessionUuid,
			Role:        userRole,
			//
			CompanyId:             companyId,
			Name:                  req.Name,
			Description:           req.Description,
			StartTime:             time.Unix(req.StartTime, 0),
			EndTime:               time.Unix(req.EndTime, 0),
			BreakDurationMinutes:  req.BreakDurationMinutes,
			GracePeriodMinutes:    req.GracePeriodMinutes,
			EarlyDepartureMinutes: req.EarlyDepartureMinutes,
			WorkDays:              req.WorkDays,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, reps)
}

// DeleteShift implements iHandler.
// @Summary      Delete shift
// @Description  Delete shift information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/shift/{id} [delete]
func (h *Handler) DeleteShift(c *gin.Context) {
	shiftId := c.Param("id")
	shiftUuid, err := uuidShared.ParseUUID(shiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift ID")
		return
	}
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call service delete shift
	errReq := applicationService.GetShiftService().DeleteShift(
		c,
		&applicationModel.DeleteShiftInput{
			// User info
			UserId:      userUuid,
			UserSession: sessionUuid,
			Role:        userRole,
			//
			ShiftId: shiftUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Delete shift successfully")
}

// DeleteShiftForUser implements iHandler.
// @Summary      Delete shift employee
// @Description  Delete shift employee information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/{id} [delete]
func (h *Handler) DeleteShiftForUser(c *gin.Context) {
	// Get shift user id
	shiftUserId := c.Param("id")
	shiftUserUuid, err := uuidShared.ParseUUID(shiftUserId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift user ID")
		return
	}
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call service delete shift for user
	errReq := applicationService.GetShiftEmployeeService().DeleteShiftUser(
		c,
		&applicationModel.DeleteShiftUserInput{
			// User info
			UserId:    userUuid,
			SessionId: sessionUuid,
			Role:      userRole,
			//
			ShiftUserId: shiftUserUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Delete shift for user successfully")
}

// DisableShiftForUser implements iHandler.
// @Summary      Disable shift for user
// @Description  Disable shift for user information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.DisableShiftForUserReq true "Disable Shift For User Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/disable [post]
func (h *Handler) DisableShiftForUser(c *gin.Context) {
	// Get req
	var req dto.DisableShiftForUserReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, 400, "Data input error")
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
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftUserUuid, err := uuidShared.ParseUUID(req.ShiftUserId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift user ID")
		return
	}
	// Call to service
	errReq := applicationService.GetShiftEmployeeService().DisableShiftUser(
		c,
		&applicationModel.DisableShiftUserInput{
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ShiftUserId: shiftUserUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Disable shift successfully")
}

// EditShift implements iHandler.
// @Summary      Edit shift information
// @Description  Edit shift information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.EditShiftReq true "Edit Shift Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/shift/edit [post]
func (h *Handler) EditShift(c *gin.Context) {
	// Get req
	var req dto.EditShiftReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, 400, "Data input error")
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
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftUuid, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift ID")
		return
	}
	companyUuid, err := uuidShared.ParseUUID(req.CompanyId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid company ID")
		return
	}
	// Call service edit shift
	errReq := applicationService.GetShiftService().EditShift(
		c,
		&applicationModel.EditShiftInput{
			UserId:      userUuid,
			UserSession: sessionUuid,
			Role:        userRole,
			//
			ShiftId:               shiftUuid,
			CompanyId:             companyUuid,
			Name:                  req.Name,
			Description:           req.Description,
			StartTime:             time.Unix(req.StartTime, 0),
			EndTime:               time.Unix(req.EndTime, 0),
			BreakDurationMinutes:  req.BreakDurationMinutes,
			GracePeriodMinutes:    req.GracePeriodMinutes,
			EarlyDepartureMinutes: req.EarlyDepartureMinutes,
			WorkDays:              req.WorkDays,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Edit shift successfully")
}

// EditShiftForUserWithEffectiveDate implements iHandler.
// @Summary      Edit shift employee effective date information
// @Description  Edit shift employee effective date information
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.ShiftEmployeeEditEffectiveDateReq true "Edit Shift Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/user/edit/effective [post]
func (h *Handler) EditShiftForUserWithEffectiveDate(c *gin.Context) {
	// Get req
	var req dto.ShiftEmployeeEditEffectiveDateReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, 400, "Data input error")
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
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call service edit shift for user with effective date
	errReq := applicationService.GetShiftEmployeeService().EditShiftForUserWithEffectiveDate(
		c,
		&applicationModel.EditShiftForUserWithEffectiveDateInput{
			UserId:    userUuid,
			SessionId: sessionUuid,
			Role:      userRole,
			//
			NewEffectiveFrom: time.Unix(req.NewEffectiveFrom, 0),
			NewEffectiveTo:   time.Unix(req.NewEffectiveTo, 0),
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Edit shift for user with effective date successfully")
}

// EnableShiftForUser implements iHandler.
// @Summary      Enable shift for user
// @Description  Enable shift for user information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.EnableShiftForUserReq true "Enable Shift For User Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/enable [post]
func (h *Handler) EnableShiftForUser(c *gin.Context) {
	var req dto.EnableShiftForUserReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, 400, "Data input error")
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
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftUserUuid, err := uuidShared.ParseUUID(req.ShiftUserId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift user ID")
		return
	}
	// Call service enable shift for user
	errReq := applicationService.GetShiftEmployeeService().EnableShiftUser(
		c,
		&applicationModel.EnableShiftUserInput{
			// User info
			UserId:    userUuid,
			SessionId: sessionUuid,
			Role:      userRole,
			//
			ShiftUserId: shiftUserUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, "Enable shift for user successfully")
}

// GetDetailShift implements iHandler.
// @Summary      Get shift detail information
// @Description  Get shift detail information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        id  path string  true  "Shift ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/shift/{id} [get]
func (h *Handler) GetDetailShift(c *gin.Context) {
	// Get id in path
	shiftId := c.Param("id")
	shiftUuid, err := uuidShared.ParseUUID(shiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift ID")
		return
	}
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call to srervice get detail shift
	reps, errReq := applicationService.GetShiftService().GetDetailShift(
		c,
		&applicationModel.GetDetailShiftInput{
			// User info
			UserId:      userUuid,
			UserSession: sessionUuid,
			Role:        userRole,
			//
			ShiftId: shiftUuid,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, reps)
}

// GetShiftForUserWithEffectiveDate implements iHandler.
// @Summary      Get shift user effective date information
// @Description  Get shift user effective date information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.ShiftEmployeeEffectiveDateReq true "Get Shift For User With Effective Date Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift [post]
func (h *Handler) GetShiftForUserWithEffectiveDate(c *gin.Context) {
	// Get req
	var req dto.ShiftEmployeeEffectiveDateReq
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, 400, "Data input error")
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
	// Get data auth from context
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Check user id
	if req.UserId != "" {
		reqUserUuid, err := uuidShared.ParseUUID(req.UserId)
		if err != nil {
			response.ErrorResponse(c, 400, "Invalid user ID")
			return
		}
		userUuid = reqUserUuid
	}
	// Call service get shift for user with effective date
	reps, errReq := applicationService.GetShiftEmployeeService().GetShiftForUserWithEffectiveDate(
		c,
		&applicationModel.GetShiftForUserWithEffectiveDateInput{
			UserId:    userUuid,
			SessionId: sessionUuid,
			Role:      userRole,
			//
			EffectiveFrom: time.Unix(req.EffectiveFrom, 0),
			EffectiveTo:   time.Unix(req.EffectiveTo, 0),
			Page:          req.Page,
			Size:          req.Size,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(c, 200, reps)
}

/**
 * New handler and impl interface
 */
func NewHandler() iHandler {
	return &Handler{}
}
