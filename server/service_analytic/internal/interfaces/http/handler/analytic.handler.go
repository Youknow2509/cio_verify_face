package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/dto"
)

// AnalyticHandler handles analytics-related HTTP requests
type AnalyticHandler struct {
	service applicationService.IAnalyticService
}

// NewAnalyticHandler creates a new analytics handler
func NewAnalyticHandler() *AnalyticHandler {
	return &AnalyticHandler{
		service: applicationService.GetAnalyticService(),
	}
}

// GetDailyReport handles GET /api/v1/reports/daily
// @Summary Get daily attendance report
// @Description Get detailed attendance report for a specific date
// @Tags Reports
// @Accept json
// @Produce json
// @Param date query string true "Report date (YYYY-MM-DD)"
// @Param company_id query string true "Company ID (UUID)"
// @Param device_id query string false "Device ID (UUID)"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /reports/daily [get]
func (h *AnalyticHandler) GetDailyReport(c *gin.Context) {
	var query dto.DailyReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid query parameters", err.Error()))
		return
	}

	// Parse date
	date, err := query.Parse()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_DATE", "Invalid date format", err.Error()))
		return
	}

	// Extract session info from context (set by HTTP middleware)
	session := getSessionFromContext(c)

	// Prepare input
	input := &applicationModel.DailyReportInput{
		Session:   session,
		Date:      date,
		CompanyID: &query.CompanyID,
	}
	if query.DeviceID != "" {
		input.DeviceID = &query.DeviceID
	}

	// Get report
	report, appErr := h.service.GetDailyReport(c.Request.Context(), input)
	if appErr != nil {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(report))
}

// GetSummaryReport handles GET /api/v1/reports/summary
// @Summary Get monthly summary report
// @Description Get monthly attendance summary with weekly breakdown
// @Tags Reports
// @Accept json
// @Produce json
// @Param month query string true "Report month (YYYY-MM)"
// @Param company_id query string true "Company ID (UUID)"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /reports/summary [get]
func (h *AnalyticHandler) GetSummaryReport(c *gin.Context) {
	var query dto.SummaryReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid query parameters", err.Error()))
		return
	}

	// Extract session info from context (set by HTTP middleware)
	session := getSessionFromContext(c)

	// Prepare input
	input := &applicationModel.SummaryReportInput{
		Session:   session,
		Month:     query.Month,
		CompanyID: &query.CompanyID,
	}

	// Get report
	report, appErr := h.service.GetSummaryReport(c.Request.Context(), input)
	if appErr != nil {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(report))
}

// ExportReport handles POST /api/v1/reports/export
// @Summary Export attendance report
// @Description Export attendance report to Excel/PDF/CSV
// @Tags Reports
// @Accept json
// @Produce json
// @Param request body dto.ExportReportRequest true "Export request"
// @Success 202 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Security Bearer
// @Router /reports/export [post]
func (h *AnalyticHandler) ExportReport(c *gin.Context) {
	var req dto.ExportReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_INPUT", "Invalid request body", err.Error()))
		return
	}

	// Extract session info from context (set by HTTP middleware)
	session := getSessionFromContext(c)

	// Prepare input
	input := &applicationModel.ExportReportInput{
		Session:   session,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Format:    req.Format,
		CompanyID: &req.CompanyID,
	}
	if req.Email != "" {
		input.Email = &req.Email
	}

	// Export report
	result, appErr := h.service.ExportReport(c.Request.Context(), input)
	if appErr != nil {
		c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}

	c.JSON(http.StatusAccepted, dto.NewSuccessResponse(result))
}

// getSessionFromContext extracts session info from gin context
// This should be set by authentication middleware
func getSessionFromContext(c *gin.Context) *applicationModel.SessionInfo {
	// Try to get session info set by middleware
	if sessionData, exists := c.Get("session"); exists {
		if session, ok := sessionData.(*applicationModel.SessionInfo); ok {
			return session
		}
	}

	// Fallback: create session from individual context values (for testing or when middleware uses different keys)
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	sessionID, _ := c.Get("session_id")
	companyID, _ := c.Get("company_id")

	session := &applicationModel.SessionInfo{}
	if uid, ok := userID.(string); ok {
		session.UserID = uid
	}
	if r, ok := role.(int32); ok {
		session.Role = r
	} else if r, ok := role.(int); ok {
		session.Role = int32(r)
	}
	if sid, ok := sessionID.(string); ok {
		session.SessionID = sid
	}
	if cid, ok := companyID.(string); ok {
		session.CompanyID = cid
	}
	session.ClientIP = c.ClientIP()
	session.ClientAgent = c.Request.UserAgent()

	return session
}
