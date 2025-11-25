package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	applicationErrors "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/errors"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/dto"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/mapper"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/middleware"
)

// ScyllaHandler handles ScyllaDB-related HTTP requests
type ScyllaHandler struct {
	service applicationService.IAnalyticService
}

// NewScyllaHandler creates a new ScyllaDB handler
func NewScyllaHandler() *ScyllaHandler {
	return &ScyllaHandler{
		service: applicationService.GetAnalyticService(),
	}
}

// ============================================
// Authorization helpers (centralized to avoid duplication)
// ============================================

// getSession extracts the *SessionInfo injected by auth middleware.
// If absent or invalid, it writes an error response and returns nil.
func getSession(c *gin.Context) *applicationModel.SessionInfo {
	sessionData, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Missing session", ""))
		return nil
	}
	session, ok := sessionData.(*applicationModel.SessionInfo)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Invalid session", ""))
		return nil
	}
	return session
}

// authorizeCompanyWide enforces: SystemAdmin full access; CompanyAdmin only own company; Employees forbidden.
// Returns session on success; on failure already responded.
func authorizeCompanyWide(c *gin.Context, companyIDStr string) *applicationModel.SessionInfo {
	session := getSession(c)
	if session == nil {
		return nil
	}
	role := domainModel.Role(session.Role)
	if role == domainModel.RoleEmployee {
		c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "Employees cannot access company-wide resources", ""))
		return nil
	}
	if role == domainModel.RoleCompanyAdmin && session.CompanyID != companyIDStr {
		c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "Cannot access another company's resources", ""))
		return nil
	}
	return session
}

// authorizeEmployeeScoped enforces: SystemAdmin any; CompanyAdmin same company any employee; Employee only self.
func authorizeEmployeeScoped(c *gin.Context, companyIDStr, employeeIDStr string) *applicationModel.SessionInfo {
	session := getSession(c)
	if session == nil {
		return nil
	}
	role := domainModel.Role(session.Role)
	if role == domainModel.RoleCompanyAdmin && session.CompanyID != companyIDStr {
		c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "Cannot access another company's employee", ""))
		return nil
	}
	if role == domainModel.RoleEmployee {
		if session.CompanyID != companyIDStr || session.UserID != employeeIDStr {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse("FORBIDDEN", "Employees can only access their own data", ""))
			return nil
		}
	}
	return session
}

// ============================================
// Attendance Records handlers
// ============================================

// GetAttendanceRecords handles GET /api/v1/attendance-records
// @Summary Get attendance records
// @Description Get attendance records for a company and month
// @Tags Attendance Records
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param limit query int false "Limit (default 100)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /attendance-records [get]
func (h *ScyllaHandler) GetAttendanceRecords(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")
	limitStr := c.DefaultQuery("limit", "100")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	records, err := h.service.GetAttendanceRecords(c.Request.Context(), companyID, yearMonth, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(records))
}

// GetAttendanceRecordsByTimeRange handles GET /api/v1/attendance-records/range
// @Summary Get attendance records by time range
// @Description Get attendance records within a time range
// @Tags Attendance Records
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param start_time query string true "Start time (RFC3339)"
// @Param end_time query string true "End time (RFC3339)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /attendance-records/range [get]
func (h *ScyllaHandler) GetAttendanceRecordsByTimeRange(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid start_time", err.Error()))
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid end_time", err.Error()))
		return
	}

	records, err := h.service.GetAttendanceRecordsByTimeRange(c.Request.Context(), companyID, yearMonth, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(records))
}

// GetAttendanceRecordsByEmployee handles GET /api/v1/attendance-records/employee/:employee_id
// @Summary Get attendance records by employee
// @Description Get attendance records for a specific employee
// @Tags Attendance Records
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID (UUID)"
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /attendance-records/employee/{employee_id} [get]
func (h *ScyllaHandler) GetAttendanceRecordsByEmployee(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")

	// Authorization
	if authorizeEmployeeScoped(c, companyIDStr, employeeIDStr) == nil {
		return
	}

	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid employee_id", err.Error()))
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	records, err := h.service.GetAttendanceRecordsByEmployee(c.Request.Context(), companyID, yearMonth, employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(records))
}

// @Description Get attendance records indexed by user
// @Tags Attendance Records
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID (UUID)"
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param limit query int false "Limit (default 100)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /attendance-records/user/{employee_id} [get]
func (h *ScyllaHandler) GetAttendanceRecordsByUser(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")
	limitStr := c.DefaultQuery("limit", "100")

	// Authorization
	if authorizeEmployeeScoped(c, companyIDStr, employeeIDStr) == nil {
		return
	}

	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid employee_id", err.Error()))
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	records, err := h.service.GetAttendanceRecordsByUser(c.Request.Context(), companyID, employeeID, yearMonth, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(records))
}

// ============================================
// Daily Summary handlers
// ============================================

// GetDailySummaries handles GET /api/v1/daily-summaries
// @Summary Get daily summaries
// @Description Get daily summaries for a month
// @Tags Daily Summaries
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param month query string true "Month (YYYY-MM)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /daily-summaries [get]
func (h *ScyllaHandler) GetDailySummaries(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	month := c.Query("month")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	summaries, err := h.service.GetDailySummaries(c.Request.Context(), companyID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summaries", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(summaries))
}

// CreateDailySummary handles POST /api/v1/daily-summaries
// @Summary Create daily summary
// @Description Create a new daily summary
// @Tags Daily Summaries
// @Accept json
// @Produce json
// @Param summary body dto.DailySummaryRequest true "Daily Summary"
// @Success 201 {object} dto.APIResponse{data=dto.DailySummaryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /daily-summaries [post]
// UpdateDailySummary handles PUT /api/v1/daily-summaries
// @Summary Update daily summary
// @Description Update an existing daily summary
// @Tags Daily Summaries
// @Accept json
// @Produce json
// @Param summary body dto.DailySummaryRequest true "Daily Summary"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /daily-summaries [put]
// GetDailySummariesByUser handles GET /api/v1/daily-summaries/user/:employee_id
// @Summary Get daily summaries by user
// @Description Get daily summaries for a specific user
// @Tags Daily Summaries
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID (UUID)"
// @Param company_id query string true "Company ID (UUID)"
// @Param month query string true "Month (YYYY-MM)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /daily-summaries/user/{employee_id} [get]
func (h *ScyllaHandler) GetDailySummariesByUser(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	companyIDStr := c.Query("company_id")
	month := c.Query("month")

	// Authorization
	if authorizeEmployeeScoped(c, companyIDStr, employeeIDStr) == nil {
		return
	}

	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid employee_id", err.Error()))
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	summaries, err := h.service.GetDailySummariesByUser(c.Request.Context(), companyID, employeeID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summaries", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(summaries))
}

// ============================================
// Audit Logs handlers
// ============================================

// GetAuditLogs handles GET /api/v1/audit-logs
// @Summary Get audit logs
// @Description Get audit logs for a company and month
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param limit query int false "Limit (default 100)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /audit-logs [get]
func (h *ScyllaHandler) GetAuditLogs(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")
	limitStr := c.DefaultQuery("limit", "100")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	logs, err := h.service.GetAuditLogs(c.Request.Context(), companyID, yearMonth, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get audit logs", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(logs))
}

// GetAuditLogsByTimeRange handles GET /api/v1/audit-logs/range
// @Summary Get audit logs by time range
// @Description Get audit logs within a time range
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param start_time query string true "Start time (RFC3339)"
// @Param end_time query string true "End time (RFC3339)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /audit-logs/range [get]
func (h *ScyllaHandler) GetAuditLogsByTimeRange(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid start_time", err.Error()))
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid end_time", err.Error()))
		return
	}

	logs, err := h.service.GetAuditLogsByTimeRange(c.Request.Context(), companyID, yearMonth, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get audit logs", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(logs))
}

// CreateAuditLog handles POST /api/v1/audit-logs
// @Summary Create audit log
// @Description Create a new audit log
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Param log body dto.AuditLogRequest true "Audit Log"
// @Success 201 {object} dto.APIResponse{data=dto.AuditLogResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /audit-logs [post]
func (h *ScyllaHandler) CreateAuditLog(c *gin.Context) {
	var req dto.AuditLogRequest

	// Validate and bind request
	if err := middleware.ValidateAndBind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid request data", err.Error()))
		return
	}

	// Map DTO to domain model
	log, err := mapper.AuditLogRequestToModel(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("MAPPING_ERROR", "Failed to map request data", err.Error()))
		return
	}

	if err := h.service.CreateAuditLog(c.Request.Context(), log); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("CREATE_FAILED", "Failed to create audit log", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse(map[string]string{"message": "Audit log created successfully"}))
}

// ============================================
// Face Enrollment Logs handlers
// ============================================

// GetFaceEnrollmentLogs handles GET /api/v1/face-enrollment-logs
// @Summary Get face enrollment logs
// @Description Get face enrollment logs for a company and month
// @Tags Face Enrollment Logs
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param limit query int false "Limit (default 100)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /face-enrollment-logs [get]
func (h *ScyllaHandler) GetFaceEnrollmentLogs(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")
	limitStr := c.DefaultQuery("limit", "100")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	logs, err := h.service.GetFaceEnrollmentLogs(c.Request.Context(), companyID, yearMonth, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get face enrollment logs", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(logs))
}

// GetFaceEnrollmentLogsByEmployee handles GET /api/v1/face-enrollment-logs/employee/:employee_id
// @Summary Get face enrollment logs by employee
// @Description Get face enrollment logs for a specific employee
// @Tags Face Enrollment Logs
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID (UUID)"
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /face-enrollment-logs/employee/{employee_id} [get]
func (h *ScyllaHandler) GetFaceEnrollmentLogsByEmployee(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")

	// Authorization
	if authorizeEmployeeScoped(c, companyIDStr, employeeIDStr) == nil {
		return
	}

	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid employee_id", err.Error()))
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	logs, err := h.service.GetFaceEnrollmentLogsByEmployee(c.Request.Context(), companyID, yearMonth, employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get face enrollment logs", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(logs))
}

// CreateFaceEnrollmentLog handles POST /api/v1/face-enrollment-logs
// @Summary Create face enrollment log
// @Description Create a new face enrollment log
// @Tags Face Enrollment Logs
// @Accept json
// @Produce json
// @Param log body dto.FaceEnrollmentLogRequest true "Face Enrollment Log"
// @Success 201 {object} dto.APIResponse{data=dto.FaceEnrollmentLogResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /face-enrollment-logs [post]
// ============================================
// Attendance Records No Shift handlers
// ============================================

// GetAttendanceRecordsNoShift handles GET /api/v1/attendance-records-no-shift
// @Summary Get attendance records without shift
// @Description Get attendance records without shift information
// @Tags Attendance Records No Shift
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param limit query int false "Limit (default 100)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /attendance-records-no-shift [get]
func (h *ScyllaHandler) GetAttendanceRecordsNoShift(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	yearMonth := c.Query("year_month")
	limitStr := c.DefaultQuery("limit", "100")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	records, err := h.service.GetAttendanceRecordsNoShift(c.Request.Context(), companyID, yearMonth, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(records))
}

// ============================================
// Company Admin - Daily Attendance Status
// ============================================

// GetDailyAttendanceStatus handles GET /api/v1/company/daily-attendance-status
// @Summary Get daily attendance status for company
// @Description Get comprehensive attendance status for a specific day including all employees
// @Tags Company Admin
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /company/daily-attendance-status [get]
func (h *ScyllaHandler) GetDailyAttendanceStatus(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	dateStr := c.Query("date")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid date format", err.Error()))
		return
	}

	yearMonth := date.Format("2006-01")
	month := date.Format("2006-01")

	// Get attendance records for the day
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC)

	records, err := h.service.GetAttendanceRecordsByTimeRange(c.Request.Context(), companyID, yearMonth, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	// Get daily summaries for the day
	summaries, err := h.service.GetDailySummaries(c.Request.Context(), companyID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summaries", err.Error()))
		return
	}

	// Filter summaries for the specific date
	var daySummaries []*domainModel.DailySummary
	for _, s := range summaries {
		if s.WorkDate.Format("2006-01-02") == dateStr {
			daySummaries = append(daySummaries, s)
		}
	}

	response := map[string]interface{}{
		"date":               date.Format("2006-01-02"),
		"total_records":      len(records),
		"total_employees":    len(daySummaries),
		"attendance_records": records,
		"daily_summaries":    daySummaries,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

// GetAttendanceStatusByTimeRange handles GET /api/v1/company/attendance-status/range
// @Summary Get attendance status for time range
// @Description Get attendance status for a company within a time range
// @Tags Company Admin
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /company/attendance-status/range [get]
func (h *ScyllaHandler) GetAttendanceStatusByTimeRange(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid start_date format", err.Error()))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid end_date format", err.Error()))
		return
	}

	// Collect all records across months in the range
	var allRecords []*domainModel.AttendanceRecord
	var allSummaries []*domainModel.DailySummary

	currentMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for currentMonth.Before(endMonth.AddDate(0, 1, 0)) {
		yearMonth := currentMonth.Format("2006-01")

		// Get attendance records
		startTime := startDate
		endTime := endDate
		if currentMonth.After(startDate) {
			startTime = currentMonth
		}
		monthEnd := currentMonth.AddDate(0, 1, -1)
		if monthEnd.Before(endDate) {
			endTime = monthEnd
		}

		records, err := h.service.GetAttendanceRecordsByTimeRange(c.Request.Context(), companyID, yearMonth, startTime, endTime)
		if err == nil {
			allRecords = append(allRecords, records...)
		}

		// Get daily summaries
		summaries, err := h.service.GetDailySummaries(c.Request.Context(), companyID, yearMonth)
		if err == nil {
			for _, s := range summaries {
				if (s.WorkDate.After(startDate) || s.WorkDate.Equal(startDate)) &&
					(s.WorkDate.Before(endDate) || s.WorkDate.Equal(endDate)) {
					allSummaries = append(allSummaries, s)
				}
			}
		}

		currentMonth = currentMonth.AddDate(0, 1, 0)
	}

	response := map[string]interface{}{
		"start_date":         startDateStr,
		"end_date":           endDateStr,
		"total_records":      len(allRecords),
		"total_summaries":    len(allSummaries),
		"attendance_records": allRecords,
		"daily_summaries":    allSummaries,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

// ============================================
// Company Admin - Monthly Detailed Summary
// ============================================

// GetMonthlyDetailedSummary handles GET /api/v1/company/monthly-summary
// @Summary Get detailed monthly summary for company
// @Description Get comprehensive monthly summary including daily summaries and records without shift
// @Tags Company Admin
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID (UUID)"
// @Param month query string true "Month (YYYY-MM)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /company/monthly-summary [get]
func (h *ScyllaHandler) GetMonthlyDetailedSummary(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	month := c.Query("month")

	// Authorization
	if authorizeCompanyWide(c, companyIDStr) == nil {
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid company_id", err.Error()))
		return
	}

	// Validate month format
	_, err = time.Parse("2006-01", month)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid month format (expected YYYY-MM)", err.Error()))
		return
	}

	// Get daily summaries for the month
	dailySummaries, err := h.service.GetDailySummaries(c.Request.Context(), companyID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summaries", err.Error()))
		return
	}

	// Get attendance records without shift
	recordsNoShift, err := h.service.GetAttendanceRecordsNoShift(c.Request.Context(), companyID, month, 1000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get records without shift", err.Error()))
		return
	}

	// Calculate statistics
	totalWorkMinutes := 0
	totalLateMinutes := 0
	totalEarlyLeaveMinutes := 0
	presentDays := 0
	lateDays := 0
	earlyLeaveDays := 0
	absentDays := 0

	employeeStats := make(map[string]map[string]interface{})

	for _, summary := range dailySummaries {
		totalWorkMinutes += summary.TotalWorkMinutes
		totalLateMinutes += summary.LateMinutes
		totalEarlyLeaveMinutes += summary.EarlyLeaveMinutes

		employeeID := summary.EmployeeID.String()
		if _, exists := employeeStats[employeeID]; !exists {
			employeeStats[employeeID] = map[string]interface{}{
				"employee_id":        employeeID,
				"present_days":       0,
				"late_days":          0,
				"early_leave_days":   0,
				"absent_days":        0,
				"total_work_minutes": 0,
			}
		}

		stats := employeeStats[employeeID]
		stats["total_work_minutes"] = stats["total_work_minutes"].(int) + summary.TotalWorkMinutes

		switch summary.AttendanceStatus {
		case 0: // PRESENT
			presentDays++
			stats["present_days"] = stats["present_days"].(int) + 1
		case 1: // LATE
			lateDays++
			presentDays++
			stats["late_days"] = stats["late_days"].(int) + 1
			stats["present_days"] = stats["present_days"].(int) + 1
		case 2: // EARLY_LEAVE
			earlyLeaveDays++
			presentDays++
			stats["early_leave_days"] = stats["early_leave_days"].(int) + 1
			stats["present_days"] = stats["present_days"].(int) + 1
		case 3: // ABSENT
			absentDays++
			stats["absent_days"] = stats["absent_days"].(int) + 1
		}
	}

	response := map[string]interface{}{
		"month":                     month,
		"total_daily_summaries":     len(dailySummaries),
		"total_records_no_shift":    len(recordsNoShift),
		"total_work_minutes":        totalWorkMinutes,
		"total_late_minutes":        totalLateMinutes,
		"total_early_leave_minutes": totalEarlyLeaveMinutes,
		"present_days":              presentDays,
		"late_days":                 lateDays,
		"early_leave_days":          earlyLeaveDays,
		"absent_days":               absentDays,
		"daily_summaries":           dailySummaries,
		"records_no_shift":          recordsNoShift,
		"employee_statistics":       employeeStats,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

// ============================================
// Company Admin - Export Endpoints
// ============================================

// ExportDailyStatus handles POST /api/v1/company/export-daily-status
// @Summary Export daily attendance status
// @Description Export daily attendance status to Excel/PDF/CSV
// @Tags Company Admin
// @Accept json
// @Produce json
// @Param request body dto.ExportDailyStatusRequest true "Export request"
// @Success 202 {object} dto.APIResponse{data=dto.ExportJobResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /company/export-daily-status [post]
func (h *ScyllaHandler) ExportDailyStatus(c *gin.Context) {
	var req dto.ExportDailyStatusRequest

	// Validate and bind request
	if err := middleware.ValidateAndBind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid request data", err.Error()))
		return
	}

	// Authorization
	if authorizeCompanyWide(c, req.CompanyID) == nil {
		return
	}

	session := authorizeCompanyWide(c, req.CompanyID)
	if session == nil {
		return
	}

	input := &applicationModel.ExportReportInput{
		Session:   session,
		StartDate: req.Date,
		EndDate:   req.Date,
		Format:    req.Format,
		CompanyID: &req.CompanyID,
	}
	if req.Email != "" {
		input.Email = &req.Email
	}

	out, aerr := h.service.ExportReport(c.Request.Context(), input)
	if aerr != nil {
		c.JSON(aerr.StatusCode, dto.NewErrorResponse(aerr.Code, aerr.Message, aerr.Details))
		return
	}
	c.JSON(http.StatusOK, dto.NewSuccessResponse(out))
}

// ExportMonthlySummary handles POST /api/v1/company/export-monthly-summary
// @Summary Export monthly detailed summary
// @Description Export monthly summary to Excel/PDF/CSV
// @Tags Company Admin
// @Accept json
// @Produce json
// @Param request body dto.ExportMonthlySummaryRequest true "Export request"
// @Success 202 {object} dto.APIResponse{data=dto.ExportJobResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /company/export-monthly-summary [post]
func (h *ScyllaHandler) ExportMonthlySummary(c *gin.Context) {
	var req dto.ExportMonthlySummaryRequest

	// Validate and bind request
	if err := middleware.ValidateAndBind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid request data", err.Error()))
		return
	}

	// Authorization
	if authorizeCompanyWide(c, req.CompanyID) == nil {
		return
	}

	session := authorizeCompanyWide(c, req.CompanyID)
	if session == nil {
		return
	}

	// Validate month format
	monthStart, err := time.Parse("2006-01", req.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(applicationErrors.ErrInvalidDateFormat.Code, "Invalid month format", err.Error()))
		return
	}
	monthEnd := monthStart.AddDate(0, 1, -1)
	startStr := monthStart.Format("2006-01-02")
	endStr := monthEnd.Format("2006-01-02")

	input := &applicationModel.ExportReportInput{
		Session:   session,
		StartDate: startStr,
		EndDate:   endStr,
		Format:    req.Format,
		CompanyID: &req.CompanyID,
	}
	if req.Email != "" {
		input.Email = &req.Email
	}

	out, aerr := h.service.ExportReport(c.Request.Context(), input)
	if aerr != nil {
		c.JSON(aerr.StatusCode, dto.NewErrorResponse(aerr.Code, aerr.Message, aerr.Details))
		return
	}
	c.JSON(http.StatusOK, dto.NewSuccessResponse(out))
}
