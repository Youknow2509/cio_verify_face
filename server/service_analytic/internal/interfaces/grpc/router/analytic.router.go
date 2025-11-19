package router

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/grpc/handler"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AnalyticRouter implements the gRPC server for analytics service
// Note: This is a placeholder until proto files are generated
// Run `bash proto-gen.sh` after installing protoc to generate actual gRPC server interface
type AnalyticRouter struct {
	handler *handler.AnalyticGrpcHandler
}

// NewAnalyticRouter creates a new analytics gRPC router
func NewAnalyticRouter() *AnalyticRouter {
	return &AnalyticRouter{
		handler: handler.NewAnalyticGrpcHandler(),
	}
}

// Placeholder methods - these will be replaced with generated gRPC methods

// GetDailyReport is a placeholder for the generated gRPC method
func (r *AnalyticRouter) GetDailyReport(ctx context.Context, req interface{}) (interface{}, error) {
	// This will be replaced with actual proto-generated implementation
	return nil, status.Errorf(codes.Unimplemented, "method GetDailyReport not implemented - run proto-gen.sh")
}

// GetSummaryReport is a placeholder for the generated gRPC method
func (r *AnalyticRouter) GetSummaryReport(ctx context.Context, req interface{}) (interface{}, error) {
	// This will be replaced with actual proto-generated implementation
	return nil, status.Errorf(codes.Unimplemented, "method GetSummaryReport not implemented - run proto-gen.sh")
}

// ExportReport is a placeholder for the generated gRPC method
func (r *AnalyticRouter) ExportReport(ctx context.Context, req interface{}) (interface{}, error) {
	// This will be replaced with actual proto-generated implementation
	return nil, status.Errorf(codes.Unimplemented, "method ExportReport not implemented - run proto-gen.sh")
}
