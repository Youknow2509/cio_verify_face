package handler

import (
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	GetListShift(*gin.Context)
	ChangeStatusShift(*gin.Context)
	// For shift employee
	GetShiftUserWithEffectiveDate(*gin.Context)
	EditShiftUserWithEffectiveDate(*gin.Context)
	EnableShiftUser(*gin.Context)
	DisableShiftUser(*gin.Context)
	DeleteShiftUser(*gin.Context)
	AddShiftEmployee(*gin.Context)
	AddShiftEmployeeList(*gin.Context)
}

/**
 * Handler struct
 */
type Handler struct{}

// AddShiftEmployeeList implements iHandler.
// @Summary      Add shift employee list
// @Description  Add shift employee list for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.AddShiftEmployeeListReq true "Add Shift Employee List Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/add/list [post]
func (h *Handler) AddShiftEmployeeList(g *gin.Context) {
	// get req
	var req dto.AddShiftEmployeeListReq
	if err := g.ShouldBind(&req); err != nil {
		response.ErrorResponse(g, 400, "Data input error")
		return
	}
	// Validate req
	validateMiddleware, ok := g.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(g, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(g, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	if err := validate.Struct(&req); err != nil {
		response.ErrorResponse(g, 400, "Validation error")
		return
	}
	// Get data auth from context
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(g)
	if !ok {
		response.ErrorResponse(g, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftUuid, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(g, 400, "Invalid shift ID")
		return
	}
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyId)
	if err != nil {
		response.ErrorResponse(g, 400, "Invalid company ID")
		return
	}
	listUserIdAdd := make([]uuid.UUID, 0)
	for _, id := range req.EmployeeIDs {
		userUuid, err := uuidShared.ParseUUID(id)
		if err != nil {
			response.ErrorResponse(g, 400, "Invalid employee ID: "+id)
			return
		}
		listUserIdAdd = append(listUserIdAdd, userUuid)
	}
	// make request to service add shift employee list
	appReq := applicationModel.AddShiftEmployeeListInput{
		// User info
		UserId:      userUuid,
		SessionId:   sessionUuid,
		Role:        userRole,
		ClientIp:    g.ClientIP(),
		ClientAgent: g.Request.UserAgent(),
		CompanyId:   companyUuid,
		//
		CompanyRequestId: companyIdReq,
		ShiftId:          shiftUuid,
		EffectiveFrom:    time.Unix(req.EffectiveFrom, 0),
		EffectiveTo:      time.Unix(req.EffectiveTo, 0),
		EmployeeIDs:      listUserIdAdd,
	}
	// Call service add shift employee list
	errReq := applicationService.GetShiftEmployeeService().AddListShiftEmployee(
		g,
		&appReq,
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(g, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(g, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(g, 200, "Add shift employee list successfully")
}

// ChangeStatusShift implements iHandler.
// @Summary      Change status shift
// @Description  Change status shift information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.ChangeStatusShiftReq true "Change Status Shift Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/shift/status [post]
func (h *Handler) ChangeStatusShift(g *gin.Context) {
	// Get req
	var req dto.ChangeStatusShiftReq
	if err := g.ShouldBind(&req); err != nil {
		response.ErrorResponse(g, 400, "Data input error")
		return
	}
	// Validate req
	validateMiddleware, ok := g.Get(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME)
	if !ok {
		response.ErrorResponse(g, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	validate, ok := validateMiddleware.(*validator.Validate)
	if !ok {
		response.ErrorResponse(g, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	if err := validate.Struct(&req); err != nil {
		response.ErrorResponse(g, 400, "Validation error")
		return
	}
	// Get data auth from context
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(g)
	if !ok {
		response.ErrorResponse(g, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftUuid, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(g, 400, "Invalid shift ID")
		return
	}
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyId)
	if err != nil {
		response.ErrorResponse(g, 400, "Invalid company ID")
		return
	}
	// Call service change status shift
	errReq := applicationService.GetShiftService().ChangeStatusShift(
		g,
		&applicationModel.ChangeStatusShiftInput{
			// User info
			UserId:      userUuid,
			UserSession: sessionUuid,
			Role:        userRole,
			ClientIp:    g.ClientIP(),
			ClientAgent: g.Request.UserAgent(),
			CompanyId:   companyUuid,
			//
			CompanyIdReq: companyIdReq,
			ShiftId:      shiftUuid,
			IsActive:     req.Status == 1,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(g, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(g, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(g, 200, "Change status shift successfully")
}

// GetListShift implements iHandler.
// @Summary      Get list shift information
// @Description  Get list shift information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        page query int false "Page number"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/shift [get]
func (h *Handler) GetListShift(g *gin.Context) {
	// Get page query
	pageStr := g.DefaultQuery("page", constants.DEFAULT_PAGE_STRING)
	// Validate and parse page
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt <= 0 {
		response.ErrorResponse(g, 400, "Invalid page number")
		return
	}
	// Get data auth from context
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(g)
	if !ok {
		response.ErrorResponse(g, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	// Call to service get list shift
	reps, errReq := applicationService.GetShiftService().GetListShift(
		g,
		&applicationModel.GetListShiftInput{
			// User info
			UserId:      userUuid,
			UserSession: sessionUuid,
			Role:        userRole,
			ClientIp:    g.ClientIP(),
			ClientAgent: g.Request.UserAgent(),
			CompanyId:   companyUuid,
			//
			Page: pageInt,
		},
	)
	if errReq != nil {
		if errReq.ErrorSystem != nil {
			response.ErrorResponse(g, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(g, 400, errReq.ErrorClient)
		return
	}
	response.SuccessResponse(g, 200, reps)
}

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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
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
	errReq := applicationService.GetShiftEmployeeService().AddShiftEmployee(
		c,
		&applicationModel.AddShiftEmployeeInput{
			// User info
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
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
	response.SuccessResponse(c, 200, "Add shift employee successfully")
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyId)
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
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
			//
			CompanyIdReq:          companyIdReq,
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
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
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
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

// DeleteShiftUser implements iHandler.
// @Summary      Delete shift employee
// @Description  Delete shift employee information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param 		 req body dto.DeleteShiftUserReq true "Delete Shift User Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/delete [post]
func (h *Handler) DeleteShiftUser(c *gin.Context) {
	// Get req
	var req dto.DeleteShiftUserReq
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	userIdReq, err := uuidShared.ParseUUID(req.EmployeeId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid employee ID")
		return
	}
	shiftId, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift ID")
		return
	}
	// Call service delete shift for user
	errReq := applicationService.GetShiftEmployeeService().DeleteShiftUser(
		c,
		&applicationModel.DeleteShiftUserInput{
			// User info
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
			//
			UserIdReq: userIdReq,
			ShiftId:   shiftId,
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

// DisableShiftUser implements iHandler.
// @Summary      Disable shift for user
// @Description  Disable shift for user information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.DisableShiftUserReq true "Disable Shift For User Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/disable [post]
func (h *Handler) DisableShiftUser(c *gin.Context) {
	// Get req
	var req dto.DisableShiftUserReq
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	companyIdUuid, _ := uuidShared.ParseUUID(companyId)
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftId, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift user ID")
		return
	}
	userIdReq, err := uuidShared.ParseUUID(req.EmployeeId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid employee ID")
		return
	}
	// Call to service
	errReq := applicationService.GetShiftEmployeeService().DisableShiftUser(
		c,
		&applicationModel.DisableShiftUserInput{
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			ShiftId:     shiftId,
			UserIdReq:   userIdReq,
			CompanyId:   companyIdUuid,
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
	}
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftUuid, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift ID")
		return
	}
	companyUuidReq, err := uuidShared.ParseUUID(req.CompanyId)
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
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
			//
			ShiftId:               shiftUuid,
			CompanyIdReq:          companyUuidReq,
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

// EditShiftUserWithEffectiveDate implements iHandler.
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
func (h *Handler) EditShiftUserWithEffectiveDate(c *gin.Context) {
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	companyUuid, _ := uuidShared.ParseUUID(companyId)
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftId, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift ID")
		return
	}
	userReqId, err := uuidShared.ParseUUID(req.EmployeeID)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid employee ID")
		return
	}
	// Call service edit shift for user with effective date
	errReq := applicationService.GetShiftEmployeeService().EditShiftForUserWithEffectiveDate(
		c,
		&applicationModel.EditShiftForUserWithEffectiveDateInput{
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
			//
			NewEffectiveFrom: time.Unix(req.NewEffectiveFrom, 0),
			NewEffectiveTo:   time.Unix(req.NewEffectiveTo, 0),
			ShiftId:          shiftId,
			UserIdReq:        userReqId,
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

// EnableShiftUser implements iHandler.
// @Summary      Enable shift for user
// @Description  Enable shift for user information for company
// @Tags         Shift
// @Accept       json
// @Produce      json
// @Param		 authorization header string true "Bearer <token>"
// @Param        dto body dto.EnableShiftUserReq true "Enable Shift For User Request"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /api/v1/employee/shift/enable [post]
func (h *Handler) EnableShiftUser(c *gin.Context) {
	var req dto.EnableShiftUserReq
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	companyUuid, _ := uuidShared.ParseUUID(companyId)
	userUuid, _ := uuidShared.ParseUUID(userId)
	sessionUuid, _ := uuidShared.ParseUUID(sessionId)
	shiftId, err := uuidShared.ParseUUID(req.ShiftId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid shift user ID")
		return
	}
	employeeId, err := uuidShared.ParseUUID(req.EmployeeId)
	if err != nil {
		response.ErrorResponse(c, 400, "Invalid employee ID")
		return
	}
	// Call service enable shift for user
	errReq := applicationService.GetShiftEmployeeService().EnableShiftUser(
		c,
		&applicationModel.EnableShiftUserInput{
			// User info
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
			//
			ShiftId:   shiftId,
			UserIdReq: employeeId,
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
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
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
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

// GetShiftUserWithEffectiveDate implements iHandler.
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
func (h *Handler) GetShiftUserWithEffectiveDate(c *gin.Context) {
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
	userId, sessionId, userRole, companyId, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeAuthSessionInvalid, "Invalid auth session")
		return
	}
	var companyUuid uuid.UUID
	if companyId != "" {
		companyUuid, _ = uuidShared.ParseUUID(companyId)
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
			UserId:      userUuid,
			SessionId:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
			CompanyId:   companyUuid,
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
