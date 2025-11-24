package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/constants"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/dto"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/middleware"
)

// EmployeeHandler handles employee self-service HTTP requests
// These endpoints are designed for employees to view their own data
type EmployeeHandler struct {
	service applicationService.IAnalyticService
}

// NewEmployeeHandler creates a new employee handler
func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{
		service: applicationService.GetAnalyticService(),
	}
}

// ============================================
// Employee Self-Service Attendance Records
// ============================================

// GetMyAttendanceRecords handles GET /api/v1/employee/my-attendance-records
// @Summary Get my attendance records
// @Description Employee views their own attendance records for a specific month
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param year_month query string false "Year-Month (YYYY-MM), defaults to current month"
// @Param limit query int false "Limit (default 100)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-attendance-records [get]
func (h *EmployeeHandler) GetMyAttendanceRecords(c *gin.Context) {
	// Get authenticated user info from context (set by auth middleware)
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	yearMonth := c.Query("year_month")
	if yearMonth == "" {
		// Default to current month if not specified
		yearMonth = time.Now().Format("2006-01")
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	// Get employee_id from session (assuming user is linked to an employee)
	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	// Get company_id from session
	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	// Use the user-optimized table for direct employee access
	records, err := h.service.GetAttendanceRecordsByUser(c.Request.Context(), companyID, employeeID, yearMonth, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(records))
}

// GetMyAttendanceRecordsInRange handles GET /api/v1/employee/my-attendance-records/range
// @Summary Get my attendance records in time range
// @Description Employee views their attendance records within a specific time range
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param year_month query string true "Year-Month (YYYY-MM)"
// @Param start_time query string true "Start time (RFC3339)"
// @Param end_time query string true "End time (RFC3339)"
// @Success 200 {object} dto.APIResponse{data=[]dto.AttendanceRecordResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-attendance-records/range [get]
func (h *EmployeeHandler) GetMyAttendanceRecordsInRange(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	yearMonth := c.Query("year_month")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if yearMonth == "" || startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Missing required parameters", ""))
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid start_time format", err.Error()))
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid end_time format", err.Error()))
		return
	}

	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	records, err := h.service.GetAttendanceRecordsByUserTimeRange(c.Request.Context(), companyID, employeeID, yearMonth, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(records))
}

// ============================================
// Employee Self-Service Daily Summaries
// ============================================

// GetMyDailySummaries handles GET /api/v1/employee/my-daily-summaries
// @Summary Get my daily summaries
// @Description Employee views their own daily attendance summaries for a month
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param month query string false "Month (YYYY-MM), defaults to current month"
// @Success 200 {object} dto.APIResponse{data=[]dto.DailySummaryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-daily-summaries [get]
func (h *EmployeeHandler) GetMyDailySummaries(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	month := c.Query("month")
	if month == "" {
		// Default to current month
		month = time.Now().Format("2006-01")
	}

	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	summaries, err := h.service.GetDailySummariesByUser(c.Request.Context(), companyID, employeeID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summaries", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(summaries))
}

// GetMyDailySummaryByDate handles GET /api/v1/employee/my-daily-summary/:date
// @Summary Get my daily summary for a specific date
// @Description Employee views their attendance summary for a specific date
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param date path string true "Date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=dto.DailySummaryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-daily-summary/{date} [get]
func (h *EmployeeHandler) GetMyDailySummaryByDate(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	dateStr := c.Param("date")
	workDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid date format (expected YYYY-MM-DD)", err.Error()))
		return
	}

	month := workDate.Format("2006-01")

	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	summary, err := h.service.GetDailySummaryByUserDate(c.Request.Context(), companyID, employeeID, month, workDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summary", err.Error()))
		return
	}

	if summary == nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "No summary found for this date", ""))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(summary))
}

// ============================================
// Employee Self-Service Statistics
// ============================================

// GetMyAttendanceStats handles GET /api/v1/employee/my-stats
// @Summary Get my attendance statistics
// @Description Employee views their attendance statistics for a month
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param month query string false "Month (YYYY-MM), defaults to current month"
// @Success 200 {object} dto.APIResponse{data=dto.EmployeeStatsResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-stats [get]
func (h *EmployeeHandler) GetMyAttendanceStats(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	// Get daily summaries to calculate statistics
	summaries, err := h.service.GetDailySummariesByUser(c.Request.Context(), companyID, employeeID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get summaries", err.Error()))
		return
	}

	// Calculate statistics
	stats := calculateEmployeeStats(summaries)

	c.JSON(http.StatusOK, dto.NewSuccessResponse(stats))
}

// ============================================
// Helper functions
// ============================================

// Helper function to get session from context
func getEmployeeSessionFromContext(c *gin.Context) *applicationModel.SessionInfo {
	// Try to get from session_info key first (used by HTTP auth middleware)
	sessionInfo, exists := c.Get(constants.ContextKeySessionInfo)
	if exists {
		if session, ok := sessionInfo.(*applicationModel.SessionInfo); ok {
			return session
		}
	}

	// Try alternative key
	session, exists := c.Get(constants.ContextKeySession)
	if exists {
		if sessionInfo, ok := session.(*applicationModel.SessionInfo); ok {
			return sessionInfo
		}
	}

	return nil
}

// Helper function to calculate employee statistics from daily summaries
func calculateEmployeeStats(summaries []*domainModel.DailySummaryByUser) map[string]interface{} {
	stats := map[string]interface{}{
		"month":                     "",
		"total_days":                len(summaries),
		"present_days":              0,
		"late_days":                 0,
		"early_leave_days":          0,
		"absent_days":               0,
		"total_work_minutes":        0,
		"total_late_minutes":        0,
		"total_early_leave_minutes": 0,
		"average_work_hours":        0.0,
		"attendance_rate":           0.0,
		"punctuality_rate":          0.0,
	}

	if len(summaries) == 0 {
		return stats
	}

	// Set month from first summary
	stats["month"] = summaries[0].SummaryMonth

	presentDays := 0
	lateDays := 0
	earlyLeaveDays := 0
	absentDays := 0
	totalWorkMinutes := 0
	totalLateMinutes := 0
	totalEarlyLeaveMinutes := 0

	for _, summary := range summaries {
		totalWorkMinutes += summary.TotalWorkMinutes
		totalLateMinutes += summary.LateMinutes
		totalEarlyLeaveMinutes += summary.EarlyLeaveMinutes

		// Attendance status: 0: PRESENT, 1: LATE, 2: EARLY_LEAVE, 3: ABSENT
		switch summary.AttendanceStatus {
		case 0: // PRESENT
			presentDays++
		case 1: // LATE
			lateDays++
			presentDays++ // Late is still present
		case 2: // EARLY_LEAVE
			earlyLeaveDays++
			presentDays++ // Early leave is still present
		case 3: // ABSENT
			absentDays++
		}
	}

	totalDays := len(summaries)
	stats["present_days"] = presentDays
	stats["late_days"] = lateDays
	stats["early_leave_days"] = earlyLeaveDays
	stats["absent_days"] = absentDays
	stats["total_work_minutes"] = totalWorkMinutes
	stats["total_late_minutes"] = totalLateMinutes
	stats["total_early_leave_minutes"] = totalEarlyLeaveMinutes

	// Calculate rates
	if totalDays > 0 {
		stats["attendance_rate"] = float64(presentDays) / float64(totalDays) * 100
		stats["average_work_hours"] = float64(totalWorkMinutes) / float64(totalDays) / 60.0

		// Punctuality rate: days without being late
		punctualDays := presentDays - lateDays
		if punctualDays < 0 {
			punctualDays = 0
		}
		stats["punctuality_rate"] = float64(punctualDays) / float64(totalDays) * 100
	}

	return stats
}

// ============================================
// Employee - Detailed Status by Day
// ============================================

// GetMyDailyStatus handles GET /api/v1/employee/my-daily-status
// @Summary Get my attendance status for a specific day
// @Description Employee views detailed attendance status for a specific day
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=dto.EmployeeDailyStatusResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-daily-status [get]
func (h *EmployeeHandler) GetMyDailyStatus(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid date format", err.Error()))
		return
	}

	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	yearMonth := date.Format("2006-01")
	month := date.Format("2006-01")

	// Get attendance records for the day
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC)

	records, err := h.service.GetAttendanceRecordsByUserTimeRange(c.Request.Context(), companyID, employeeID, yearMonth, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	// Get daily summary for the day
	summary, err := h.service.GetDailySummaryByUserDate(c.Request.Context(), companyID, employeeID, month, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summary", err.Error()))
		return
	}

	response := map[string]interface{}{
		"date":               dateStr,
		"total_records":      len(records),
		"attendance_records": records,
		"daily_summary":      summary,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

// ============================================
// Employee - Status by Time Range
// ============================================

// GetMyStatusByTimeRange handles GET /api/v1/employee/my-status/range
// @Summary Get my attendance status for a time range
// @Description Employee views attendance status within a specific time range
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=dto.EmployeeStatusRangeResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-status/range [get]
func (h *EmployeeHandler) GetMyStatusByTimeRange(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

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

	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	// Collect all records and summaries across months in the range
	var allRecords []*domainModel.AttendanceRecordByUser
	var allSummaries []*domainModel.DailySummaryByUser

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

		records, err := h.service.GetAttendanceRecordsByUserTimeRange(c.Request.Context(), companyID, employeeID, yearMonth, startTime, endTime)
		if err == nil {
			allRecords = append(allRecords, records...)
		}

		// Get daily summaries
		summaries, err := h.service.GetDailySummariesByUser(c.Request.Context(), companyID, employeeID, yearMonth)
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
// Employee - Detailed Monthly Summary
// ============================================

// GetMyDetailedMonthlySummary handles GET /api/v1/employee/my-monthly-summary
// @Summary Get my detailed monthly summary
// @Description Employee views comprehensive monthly attendance summary
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param month query string false "Month (YYYY-MM), defaults to current month"
// @Success 200 {object} dto.APIResponse{data=dto.EmployeeMonthlySummaryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/my-monthly-summary [get]
func (h *EmployeeHandler) GetMyDetailedMonthlySummary(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	// Validate month format
	_, err := time.Parse("2006-01", month)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid month format (expected YYYY-MM)", err.Error()))
		return
	}

	employeeID, err := uuid.Parse(session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_USER", "Invalid user ID", err.Error()))
		return
	}

	companyID, err := uuid.Parse(session.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_COMPANY", "Invalid company ID", err.Error()))
		return
	}

	// Get daily summaries for the month
	dailySummaries, err := h.service.GetDailySummariesByUser(c.Request.Context(), companyID, employeeID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get daily summaries", err.Error()))
		return
	}

	// Get attendance records for the month
	records, err := h.service.GetAttendanceRecordsByUser(c.Request.Context(), companyID, employeeID, month, 1000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("QUERY_FAILED", "Failed to get attendance records", err.Error()))
		return
	}

	// Calculate statistics
	stats := calculateEmployeeStats(dailySummaries)

	response := map[string]interface{}{
		"month":                 month,
		"total_daily_summaries": len(dailySummaries),
		"total_records":         len(records),
		"statistics":            stats,
		"daily_summaries":       dailySummaries,
		"attendance_records":    records,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

// ============================================
// Employee - Export Endpoints
// ============================================

// ExportMyDailyStatus handles POST /api/v1/employee/export-daily-status
// @Summary Export my daily attendance status
// @Description Export employee's daily attendance status to Excel/PDF/CSV
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param request body dto.ExportEmployeeDailyStatusRequest true "Export request"
// @Success 202 {object} dto.APIResponse{data=dto.ExportJobResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/export-daily-status [post]
func (h *EmployeeHandler) ExportMyDailyStatus(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	var req dto.ExportEmployeeDailyStatusRequest

	// Validate and bind request
	if err := middleware.ValidateAndBind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid request data", err.Error()))
		return
	}

	// Parse and validate date
	_, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid date format (expected YYYY-MM-DD)", err.Error()))
		return
	}

	// Prepare export input for single day (employee's own data)
	companyID := session.CompanyID
	input := &applicationModel.ExportReportInput{
		Session:   session,
		StartDate: req.Date,
		EndDate:   req.Date,
		Format:    req.Format,
		CompanyID: &companyID,
	}
	if req.Email != "" {
		input.Email = &req.Email
	}

	// Call export service (it will automatically filter to employee's own data based on role)
	result, appErr := h.service.ExportReport(c.Request.Context(), input)
	if appErr != nil {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}

	c.JSON(http.StatusAccepted, dto.NewSuccessResponse(result))
}

// ExportMyMonthlySummary handles POST /api/v1/employee/export-monthly-summary
// @Summary Export my monthly summary
// @Description Export employee's monthly summary to Excel/PDF/CSV
// @Tags Employee Self-Service
// @Accept json
// @Produce json
// @Param request body dto.ExportEmployeeMonthlySummaryRequest true "Export request"
// @Success 202 {object} dto.APIResponse{data=dto.ExportJobResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /employee/export-monthly-summary [post]
func (h *EmployeeHandler) ExportMyMonthlySummary(c *gin.Context) {
	session := getEmployeeSessionFromContext(c)
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Session not found", ""))
		return
	}

	var req dto.ExportEmployeeMonthlySummaryRequest

	// Validate and bind request
	if err := middleware.ValidateAndBind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid request data", err.Error()))
		return
	}

	// Parse and validate month
	monthStart, err := time.Parse("2006-01", req.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid month format (expected YYYY-MM)", err.Error()))
		return
	}

	// Calculate month range (first day to last day)
	monthEnd := monthStart.AddDate(0, 1, -1)
	startDate := monthStart.Format("2006-01-02")
	endDate := monthEnd.Format("2006-01-02")

	// Prepare export input for full month (employee's own data)
	companyID := session.CompanyID
	input := &applicationModel.ExportReportInput{
		Session:   session,
		StartDate: startDate,
		EndDate:   endDate,
		Format:    req.Format,
		CompanyID: &companyID,
	}
	if req.Email != "" {
		input.Email = &req.Email
	}

	// Call export service (it will automatically filter to employee's own data based on role)
	result, appErr := h.service.ExportReport(c.Request.Context(), input)
	if appErr != nil {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}

	c.JSON(http.StatusAccepted, dto.NewSuccessResponse(result))
}
