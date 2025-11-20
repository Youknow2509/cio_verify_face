package handler

import (
	"strings"
	"time"

	gin "github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/constants"
	dto "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/dto"
	response "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/response"
	contextShared "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/context"
	uuidShared "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/uuid"
)

// ============================================
// Attendance handler
// ============================================
type iAttendanceHandler interface {
	AddAttendance(c *gin.Context)
	GetAttendanceRecords(c *gin.Context)
	GetAttendanceRecordsEmployee(c *gin.Context)
	GetDailyAttendanceSummary(c *gin.Context)
	GetDailyAttendanceSummaryEmployee(c *gin.Context)
}

// ============================================================
// Attendance handler struct deployment interface
// ============================================================
type AttendanceHandler struct{}

// GetDailyAttendanceSummary implements iAttendanceHandler.
// @Summary      Get daily attendance summary
// @Description  Get daily attendance summary
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param 	  	 Authorization header string true "With the bearer started"
// @Param        request   body dto.GetDailyAttendanceSummaryRequest  true  "Request body get daily attendance summary"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/records/summary/daily [post]
func (a *AttendanceHandler) GetDailyAttendanceSummary(c *gin.Context) {
	var req *dto.GetDailyAttendanceSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(req); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		response.BadRequestResponse(
			c,
			response.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
		return
	}
	// Get session
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
	// Create session req
	sessionReq := applicationModel.SessionReq{
		UserId:      userUuid,
		SessionId:   sessionUuid,
		Role:        userRole,
		CompanyId:   companyUuid,
		ClientIp:    c.ClientIP(),
		ClientAgent: c.Request.UserAgent(),
	}
	// Parse data request
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid company_id")
		return
	}
	summaryMonth := req.SummaryMonth
	if len(summaryMonth) != 7 || summaryMonth[4] != '-' {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid summary_month format, expected YYYY-MM")
		return
	}
	var workDate time.Time
	if req.WorkDate != 0 {
		workDate = time.Unix(req.WorkDate, 0)
	}
	pageStageByte := []byte(summaryMonth)
	// Call application service
	summary, errApplication := applicationService.GetAttendanceService().GetDailyAttendanceSummaryForCompany(
		c,
		&applicationModel.GetDailyAttendanceSummaryModel{
			Session: &sessionReq,
			//
			CompanyID:    companyIdReq,
			SummaryMonth: summaryMonth,
			WorkDate:     workDate,
			PageSize:     req.PageSize,
			PageStage:    pageStageByte,
		},
	)
	if errApplication != nil {
		if errApplication.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Server temporary busy, please try again later")
			return
		}
		response.BadRequestResponse(c, 400, errApplication.ErrorClient)
		return
	}
	// Return response
	response.SuccessResponse(c, 200, summary)
}

// GetDailyAttendanceSummaryEmployee implements iAttendanceHandler.
// @Summary      Get daily attendance summary for employee
// @Description  Get daily attendance summary for employee
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param 	  	 Authorization header string true "With the bearer started"
// @Param        request   body dto.GetDailyAttendanceSummaryEmployeeRequest  true  "Request body get daily attendance summary for employee"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/records/employee/summary/daily  [post]
func (a *AttendanceHandler) GetDailyAttendanceSummaryEmployee(c *gin.Context) {
	var req *dto.GetDailyAttendanceSummaryEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(req); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		response.BadRequestResponse(
			c,
			response.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
		return
	}
	// Get session
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
	// Create session req
	sessionReq := applicationModel.SessionReq{
		UserId:      userUuid,
		SessionId:   sessionUuid,
		Role:        userRole,
		CompanyId:   companyUuid,
		ClientIp:    c.ClientIP(),
		ClientAgent: c.Request.UserAgent(),
	}
	// Parse data request
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid company_id")
		return
	}
	employeeIdReq, err := uuidShared.ParseUUID(req.EmployeeID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid employee_id")
		return
	}
	summaryMonth := req.SummaryMonth
	if len(summaryMonth) != 7 || summaryMonth[4] != '-' {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid summary_month format, expected YYYY-MM")
		return
	}
	pageStageByte := []byte(req.PageStage)
	// Call application service
	summary, errApplication := applicationService.GetAttendanceService().GetDailyAttendanceSummaryEmployeeForCompany(
		c,
		&applicationModel.GetDailyAttendanceSummaryEmployeeModel{
			Session: &sessionReq,
			//
			CompanyID:    companyIdReq,
			SummaryMonth: summaryMonth,
			EmployeeID:   employeeIdReq,
			PageSize:     req.PageSize,
			PageStage:    pageStageByte,
		},
	)
	if errApplication != nil {
		if errApplication.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Server temporary busy, please try again later")
			return
		}
		response.BadRequestResponse(c, 400, errApplication.ErrorClient)
		return
	}
	// Return response
	response.SuccessResponse(c, 200, summary)
}

// GetAttendanceRecordsEmployee implements iAttendanceHandler.
// @Summary      Get attendance records for employee
// @Description  Get attendance records for employee
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param 	  	 Authorization header string true "With the bearer started"
// @Param        request   body dto.GetAttendanceRecordsEmployeeRequest  true  "Request body get attendance records for employee"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/records/employee [post]
func (a *AttendanceHandler) GetAttendanceRecordsEmployee(c *gin.Context) {
	var req *dto.GetAttendanceRecordsEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(req); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		response.BadRequestResponse(
			c,
			response.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
		return
	}
	// Get session
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
	// Create session req
	sessionReq := applicationModel.SessionReq{
		UserId:      userUuid,
		SessionId:   sessionUuid,
		Role:        userRole,
		CompanyId:   companyUuid,
		ClientIp:    c.ClientIP(),
		ClientAgent: c.Request.UserAgent(),
	}
	// Parse date request
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid company_id")
		return
	}
	var employeeIdReq uuid.UUID
	if req.EmployeeID != "" {
		employeeIdReq, err = uuidShared.ParseUUID(req.EmployeeID)
		if err != nil {
			response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid employee_id")
			return
		}
	}
	yearMonth := req.YearMonth
	if len(yearMonth) != 7 || yearMonth[4] != '-' {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid year_month format, expected YYYY-MM")
		return
	}
	pageStageByte := []byte(yearMonth)
	// Call application service
	records, errApplication := applicationService.GetAttendanceService().GetAttendanceRecordsEmployeeForConpany(
		c,
		&applicationModel.GetAttendanceRecordsEmployeeModel{
			Session: &sessionReq,
			//
			CompanyID:  companyIdReq,
			YearMonth:  yearMonth,
			EmployeeID: employeeIdReq,
			PageSize:   req.PageSize,
			PageStage:  pageStageByte,
		},
	)
	if errApplication != nil {
		if errApplication.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Server temporary busy, please try again later")
			return
		}
		response.BadRequestResponse(c, 400, errApplication.ErrorClient)
		return
	}
	// Return response
	response.SuccessResponse(c, 200, records)
}

// GetAttendanceRecords implements iAttendanceHandler.
// @Summary      Get attendance records
// @Description  Get attendance records
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param 	  	 Authorization header string true "With the bearer started"
// @Param        request   body dto.GetAttendanceRecordsRequest  true  "Request body get attendance records"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance/records [post]
func (a *AttendanceHandler) GetAttendanceRecords(c *gin.Context) {
	var req *dto.GetAttendanceRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(req); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		response.BadRequestResponse(
			c,
			response.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
		return
	}
	// Get session
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
	// Create session req
	sessionReq := applicationModel.SessionReq{
		UserId:      userUuid,
		SessionId:   sessionUuid,
		Role:        userRole,
		CompanyId:   companyUuid,
		ClientIp:    c.ClientIP(),
		ClientAgent: c.Request.UserAgent(),
	}
	// Parse date request
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid company_id")
		return
	}
	yearMonth := req.YearMonth
	if len(yearMonth) != 7 || yearMonth[4] != '-' {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid year_month format, expected YYYY-MM")
		return
	}
	pageStageByte := []byte(yearMonth)
	// Call application service
	records, errApplication := applicationService.GetAttendanceService().GetAttendanceRecordsCompany(
		c,
		&applicationModel.GetAttendanceRecordsCompanyModel{
			Session: &sessionReq,
			//
			CompanyID: companyIdReq,
			YearMonth: yearMonth,
			PageSize:  req.PageSize,
			PageStage: pageStageByte,
		},
	)
	if errApplication != nil {
		if errApplication.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Server temporary busy, please try again later")
			return
		}
		response.BadRequestResponse(c, 400, errApplication.ErrorClient)
		return
	}
	// Return response
	response.SuccessResponse(c, 200, records)
}

// AddAttendance implements iAttendanceHandler.
// @Summary      Add attendance record
// @Description  Add attendance record
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param 	  	 Authorization header string true "With the bearer started"
// @Param        request   body dto.AddAttendanceRequest  true  "Request body add attendance record"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/attendance [post]
func (a *AttendanceHandler) AddAttendance(c *gin.Context) {
	var req *dto.AddAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(req); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		response.BadRequestResponse(
			c,
			response.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
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
	// Create session req
	sessionReq := applicationModel.SessionReq{
		UserId:      userUuid,
		SessionId:   sessionUuid,
		Role:        userRole,
		CompanyId:   companyUuid,
		ClientIp:    c.ClientIP(),
		ClientAgent: c.Request.UserAgent(),
	}
	// Parse date request
	companyIdReq, err := uuidShared.ParseUUID(req.CompanyID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid company_id")
		return
	}
	employeeIdReq, err := uuidShared.ParseUUID(req.EmployeeID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid employee_id")
		return
	}
	deviceIdReq, err := uuidShared.ParseUUID(req.DeviceID)
	if err != nil {
		response.BadRequestResponse(c, response.ErrCodeParamInvalid, "Invalid device_id")
		return
	}
	recordTime := time.Unix(req.RecordTime, 0)
	// Call application service
	errApplication := applicationService.GetAttendanceService().AddAttendance(
		c,
		&applicationModel.AddAttendanceModel{
			Session: &sessionReq,
			// Map request to model
			CompanyID:           companyIdReq,
			EmployeeID:          employeeIdReq,
			DeviceID:            deviceIdReq,
			RecordTime:          recordTime,
			VerificationMethod:  req.VerificationMethod,
			VerificationScore:   req.VerificationScore,
			FaceImageURL:        req.FaceImageURL,
			LocationCoordinates: req.LocationCoordinates,
		},
	)
	if errApplication != nil {
		if errApplication.ErrorSystem != nil {
			response.ErrorResponse(c, response.ErrorCodeSystemTemporary, "Server temporary busy, please try again later")
			return
		}
		response.BadRequestResponse(c, 400, errApplication.ErrorClient)
		return
	}
	// Return response
	response.SuccessResponse(c, 200, "Attendance record added successfully")
}

// NewAttendanceHandler creates a new instance of AttendanceHandler
func NewAttendanceHandler() iAttendanceHandler {
	return &AttendanceHandler{}
}
