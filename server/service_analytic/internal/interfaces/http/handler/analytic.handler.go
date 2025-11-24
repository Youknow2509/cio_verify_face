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
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Invalid or missing session information", ""))
		return
	}

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
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Invalid or missing session information", ""))
		return
	}

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
	if session == nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("UNAUTHORIZED", "Invalid or missing session information", ""))
		return
	}

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

// getSessionFromContext extracts session info from gin context.
// It relies *only* on the "session" object set by the authentication middleware.
// This is a more secure implementation that avoids unsafe fallbacks.
func getSessionFromContext(c *gin.Context) *applicationModel.SessionInfo {
	sessionData, exists := c.Get("session")
	if !exists {
		return nil // If "session" is not in context, there is no valid session.
	}

	session, ok := sessionData.(*applicationModel.SessionInfo)
	if !ok {
		return nil // If the data is not the correct type, the session is invalid.
	}

	// Add client IP and User Agent for logging and security purposes.
	session.ClientIP = c.ClientIP()
	session.ClientAgent = c.Request.UserAgent()

	return session
}
