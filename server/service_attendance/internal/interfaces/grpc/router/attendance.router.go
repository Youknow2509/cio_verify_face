package routes

import (
	"context"

	pb "github.com/youknow2509/cio_verify_face/server/service_attendance/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ================================
// Attendance grpc routes
// ================================
type AttendanceRouter struct {
	pb.UnimplementedAttendanceServiceServer
}

func (a *AttendanceRouter) CheckInAttendance(ctx context.Context, req *pb.CheckIOAttendanceRequest) (*pb.ResponseBase, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckInAttendance not implemented")
}

func (a *AttendanceRouter) CheckOutAttendance(ctx context.Context, req *pb.CheckIOAttendanceRequest) (*pb.ResponseBase, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckOutAttendance not implemented")
}

func (a *AttendanceRouter) GetAttendanceRecords(ctx context.Context, req *pb.GetAttendanceRecordsRequest) (*pb.GetAttendanceRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAttendanceRecords not implemented")
}

// new attendance router
func NewAttendanceRouter() *AttendanceRouter {
	return &AttendanceRouter{}
}