package handler

import (
	"context"
	"time"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AnalyticGrpcHandler handles analytics-related gRPC requests
type AnalyticGrpcHandler struct {
	service applicationService.IAnalyticService
}

// NewAnalyticGrpcHandler creates a new analytics gRPC handler
func NewAnalyticGrpcHandler() *AnalyticGrpcHandler {
	return &AnalyticGrpcHandler{
		service: applicationService.GetAnalyticService(),
	}
}

// GetDailyReport handles gRPC GetDailyReport request
// gRPC receives already-validated session info from inter-service calls
func (h *AnalyticGrpcHandler) GetDailyReport(ctx context.Context, date, companyID, locationID string) (interface{}, error) {
	// Extract session info from context (set by SessionInterceptor middleware)
	sessionInfo, err := middleware.GetSessionInfoFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "session info not found: %v", err)
	}

	global.Logger.Info("gRPC GetDailyReport request",
		"user_id", sessionInfo.UserID,
		"company_id", companyID,
		"date", date)

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
	}

	// Prepare input with session info
	input := &applicationModel.DailyReportInput{
		Session:   sessionInfo,
		Date:      parsedDate,
		CompanyID: &companyID,
	}

	// Get report (authorization handled in application service)
	report, appErr := h.service.GetDailyReport(ctx, input)
	if appErr != nil {
		return nil, status.Errorf(codes.Code(appErr.StatusCode/100), "%s: %s", appErr.Message, appErr.Details)
	}

	return report, nil
}

// GetSummaryReport handles gRPC GetSummaryReport request
// gRPC receives already-validated session info from inter-service calls
func (h *AnalyticGrpcHandler) GetSummaryReport(ctx context.Context, month, companyID string) (interface{}, error) {
	// Extract session info from context (set by SessionInterceptor middleware)
	sessionInfo, err := middleware.GetSessionInfoFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "session info not found: %v", err)
	}

	global.Logger.Info("gRPC GetSummaryReport request",
		"user_id", sessionInfo.UserID,
		"company_id", companyID,
		"month", month)

	// Prepare input with session info
	input := &applicationModel.SummaryReportInput{
		Session:   sessionInfo,
		Month:     month,
		CompanyID: &companyID,
	}

	// Get report (authorization handled in application service)
	report, appErr := h.service.GetSummaryReport(ctx, input)
	if appErr != nil {
		return nil, status.Errorf(codes.Code(appErr.StatusCode/100), "%s: %s", appErr.Message, appErr.Details)
	}

	return report, nil
}

// ExportReport handles gRPC ExportReport request
// gRPC receives already-validated session info from inter-service calls
func (h *AnalyticGrpcHandler) ExportReport(ctx context.Context, startDate, endDate, format, companyID, email string) (interface{}, error) {
	// Extract session info from context (set by SessionInterceptor middleware)
	sessionInfo, err := middleware.GetSessionInfoFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "session info not found: %v", err)
	}

	global.Logger.Info("gRPC ExportReport request",
		"user_id", sessionInfo.UserID,
		"company_id", companyID,
		"format", format)

	// Prepare input with session info
	input := &applicationModel.ExportReportInput{
		Session:   sessionInfo,
		StartDate: startDate,
		EndDate:   endDate,
		Format:    format,
		CompanyID: &companyID,
	}
	if email != "" {
		input.Email = &email
	}

	// Export report (authorization handled in application service)
	result, appErr := h.service.ExportReport(ctx, input)
	if appErr != nil {
		return nil, status.Errorf(codes.Code(appErr.StatusCode/100), "%s: %s", appErr.Message, appErr.Details)
	}

	return result, nil
}
