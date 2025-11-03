package handler

import (
	"time"

	gin "github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	appModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	appService "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/constants"
	dto "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/dto"
	response "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/response"
	contextShared "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/context"
	sharedUuid "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/uuid"
)

// ============================================
// Attendance handler
// ============================================
type iAttendanceHandler interface {
	// CheckIn handles check-in requests
	CheckIn(c *gin.Context)
	// CheckOut handles check-out requests
	CheckOut(c *gin.Context)
	// GetRecords retrieves attendance records
	GetRecords(c *gin.Context)
	// GetRecordByID retrieves a specific attendance record by ID
	GetRecordByID(c *gin.Context)
	// GetMyHistory retrieves the attendance history for the current user
	GetMyHistory(c *gin.Context)
}

// ============================================================
// Attendance handler struct deployment interface
// ============================================================
type AttendanceHandler struct{}

// CheckIn implements iAttendanceHandler.
// @Summary      Check In attendance
// @Description  User check-in for attendance
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param        authorization header string true "Bearer token"
// @Param        request   body dto.CheckInRequest  true  "Request body check-in"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/check_in [post]
func (a *AttendanceHandler) CheckIn(c *gin.Context) {
	// Get req
	var req dto.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestResponse(
			c,
			400,
			"Data request error",
		)
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
	// Get session
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Validate uuid
	userUuid, _ := sharedUuid.ParseUUID(userId)
	sessionUuid, _ := sharedUuid.ParseUUID(sessionId)
	deviceUuid, err := sharedUuid.ParseUUID(req.DeviceId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Device ID is not valid UUID")
		return
	}
	// Call service to process check-in
	errReps := appService.GetAttendanceService().CheckInUser(
		c,
		&appModel.CheckInInput{
			Timestamp: req.Timestamp,
			DeviceId:  deviceUuid,
			Location:  req.Location,
			// Session info
			UserID:      userUuid,
			SessionID:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
		},
	)
	if errReps != nil {
		if errReps.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReps.ErrorClient)
		return
	}
	// Respond success
	response.SuccessResponse(
		c,
		200,
		"Check-in successful",
	)
}

// CheckOut implements iAttendanceHandler.
// @Summary      Check Out attendance
// @Description  User check-out for attendance
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param        authorization header string true "Bearer token"
// @Param        request   body dto.CheckOutRequest  true  "Request body check-out"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/check_out [post]
func (a *AttendanceHandler) CheckOut(c *gin.Context) {
	// Get req
	var req dto.CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestResponse(
			c,
			400,
			"Data request error",
		)
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
	// Get session
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Validate uuid
	userUuid, _ := sharedUuid.ParseUUID(userId)
	sessionUuid, _ := sharedUuid.ParseUUID(sessionId)
	deviceUuid, err := sharedUuid.ParseUUID(req.DeviceId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Device ID is not valid UUID")
		return
	}
	// Call service to process check-out
	errReps := appService.GetAttendanceService().CheckOutUser(
		c,
		&appModel.CheckOutInput{
			Timestamp: req.Timestamp,
			DeviceId:  deviceUuid,
			Location:  req.Location,
			// Session info
			UserID:      userUuid,
			SessionID:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
		},
	)
	if errReps != nil {
		if errReps.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReps.ErrorClient)
		return
	}
	// Respond success
	response.SuccessResponse(
		c,
		200,
		"Check-out successful",
	)
}

// GetMyHistory implements iAttendanceHandler.
// @Summary      Get My Attendance History
// @Description  Retrieve the attendance history for the current user
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param        authorization header string true "Bearer token"
// @Param        req  body dto.GetMyAttendanceRecordsRequest  true  "Request body to get my attendance records"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/history/my [post]
func (a *AttendanceHandler) GetMyHistory(c *gin.Context) {
	// Get req
	var req dto.GetMyAttendanceRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestResponse(
			c,
			400,
			"Data request error",
		)
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
	// Get session
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Validate uuid
	userUuid, _ := sharedUuid.ParseUUID(userId)
	sessionUuid, _ := sharedUuid.ParseUUID(sessionId)
	// Validate date format
	var (
		startDate time.Time
		endDate   time.Time
	)
	if req.StartDate != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Start date is not valid format YYYY-MM-DD")
			return
		}
	}
	if req.EndDate != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.ErrorResponse(c, response.ErrorCodeValidateRequest, "End date is not valid format YYYY-MM-DD")
			return
		}
	}
	if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(startDate) {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "End date must be after start date")
		return
	}
	// Set default pagination if not provided
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	// Call service to get my records
	records, errReps := appService.GetAttendanceService().GetMyRecords(
		c,
		&appModel.GetMyRecordsInput{
			Page:      req.Page,
			Size:      req.PageSize,
			StartDate: startDate,
			EndDate:   endDate,
			// Session info
			UserID:    userUuid,
			SessionID: sessionUuid,
			Role:      userRole,
			ClientIp:  c.ClientIP(),
		},
	)
	if errReps != nil {
		if errReps.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReps.ErrorClient)
		return
	}
	// Respond success
	response.SuccessResponse(
		c,
		200,
		records,
	)
}

// GetRecordByID implements iAttendanceHandler.
// @Summary      Get Attendance Record by ID
// @Description  Retrieve a specific attendance record by ID
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param        authorization header string true "Bearer token"
// @Param        record_id  path  string  true  "Record ID"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/records/{record_id} [get]
func (a *AttendanceHandler) GetRecordByID(c *gin.Context) {
	// get record_id from path
	recordId := c.Param("record_id")
	if recordId == "" {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Record ID is required")
		return
	}
	// Get session
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Validate uuid
	userUuid, _ := sharedUuid.ParseUUID(userId)
	sessionUuid, _ := sharedUuid.ParseUUID(sessionId)
	recordUuid, err := sharedUuid.ParseUUID(recordId)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Record ID is not valid UUID")
		return
	}
	// Call service to get record by ID
	reps, errReps := appService.GetAttendanceService().GetRecordByID(
		c,
		&appModel.GetAttendanceRecordByIDInput{
			RecordID: recordUuid,
			// Session info
			UserID:      userUuid,
			SessionID:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
		},
	)
	if errReps != nil {
		if errReps.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReps.ErrorClient)
		return
	}
	// Respond success
	response.SuccessResponse(
		c,
		200,
		reps,
	)
}

// GetRecords implements iAttendanceHandler.
// @Summary      Get Attendance Records
// @Description  Retrieve attendance records for device, day, user, ...
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param        authorization header string true "Bearer token"
// @Param        req  body dto.GetAttendanceRecordsRequest  true  "Request body to get attendance records"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/records [post]
func (a *AttendanceHandler) GetRecords(c *gin.Context) {
	// Get req
	var req dto.GetAttendanceRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestResponse(
			c,
			400,
			"Data request error",
		)
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
	// Get session
	userId, sessionId, userRole, ok := contextShared.GetSessionFromContext(c)
	if !ok {
		response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
		return
	}
	// Validate uuid
	userUuid, _ := sharedUuid.ParseUUID(userId)
	sessionUuid, _ := sharedUuid.ParseUUID(sessionId)
	var deviceUuid uuid.UUID
	if req.DeviceId != "" {
		deviceUuid, err = sharedUuid.ParseUUID(req.DeviceId)
		if err != nil {
			response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Device ID is not valid UUID")
			return
		}
	}
	// Validate date format
	var (
		startDate time.Time
		endDate   time.Time
	)
	if req.StartDate != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.ErrorResponse(c, response.ErrorCodeValidateRequest, "Start date is not valid format YYYY-MM-DD")
			return
		}
	}
	if req.EndDate != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.ErrorResponse(c, response.ErrorCodeValidateRequest, "End date is not valid format YYYY-MM-DD")
			return
		}
	}
	if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(startDate) {
		response.ErrorResponse(c, response.ErrorCodeValidateRequest, "End date must be after start date")
		return
	}
	// Set default pagination if not provided
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	// Call service to get records
	records, errReps := appService.GetAttendanceService().GetRecords(
		c,
		&appModel.GetAttendanceRecordsInput{
			Page:      req.Page,
			Size:      req.PageSize,
			StartDate: startDate,
			EndDate:   endDate,
			DeviceID:  deviceUuid,
			UserID:    userUuid,
			// Session info
			SessionID:   sessionUuid,
			Role:        userRole,
			ClientIp:    c.ClientIP(),
			ClientAgent: c.Request.UserAgent(),
		},
	)
	if errReps != nil {
		if errReps.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Internal server error")
			return
		}
		response.ErrorResponse(c, 400, errReps.ErrorClient)
		return
	}
	// Respond success
	response.SuccessResponse(
		c,
		200,
		records,
	)
}

// NewAttendanceHandler creates a new instance of AttendanceHandler
func NewAttendanceHandler() iAttendanceHandler {
	return &AttendanceHandler{}
}
