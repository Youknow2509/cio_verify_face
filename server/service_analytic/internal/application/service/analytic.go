package service

import (
	"context"
	"errors"

	applicationErrors "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/errors"
	model "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
)

// IAnalyticService interface defines the analytics service operations
type IAnalyticService interface {
	// GetDailyReport returns daily attendance report
	GetDailyReport(ctx context.Context, input *model.DailyReportInput) (*model.DailyReportOutput, *applicationErrors.Error)
	
	// GetSummaryReport returns monthly summary report
	GetSummaryReport(ctx context.Context, input *model.SummaryReportInput) (*model.SummaryReportOutput, *applicationErrors.Error)
	
	// ExportReport exports attendance report to file
	ExportReport(ctx context.Context, input *model.ExportReportInput) (*model.ExportReportOutput, *applicationErrors.Error)
	
	// GetHealthCheck returns service health status
	GetHealthCheck(ctx context.Context) (*model.HealthCheckOutput, *applicationErrors.Error)
}

// Manager instance of analytic service
var _vIAnalyticService IAnalyticService

// GetAnalyticService returns the singleton instance
func GetAnalyticService() IAnalyticService {
	return _vIAnalyticService
}

// SetAnalyticService sets the singleton instance
func SetAnalyticService(service IAnalyticService) error {
	if service == nil {
		return errors.New("analytic service set is nil")
	}
	if _vIAnalyticService != nil {
		return errors.New("analytic service is already set")
	}
	_vIAnalyticService = service
	return nil
}
