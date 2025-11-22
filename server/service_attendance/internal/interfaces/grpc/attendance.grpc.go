package grpc

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
	pb "github.com/youknow2509/cio_verify_face/server/service_attendance/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AttendanceGRPCServer struct {
	pb.UnimplementedAttendanceServiceServer
	attendanceService service.IAttendanceService
}

func NewAttendanceGRPCServer() *AttendanceGRPCServer {
	attendanceService := service.GetAttendanceService()
	return &AttendanceGRPCServer{
		attendanceService: attendanceService,
	}
}

func (s *AttendanceGRPCServer) DeleteAttendanceRecords(ctx context.Context, req *pb.DeleteAttendanceRecordsInput) (*emptypb.Empty, error) {
	// Parse request to application model
	var sessionUser applicationModel.SessionReq
	var sessionService applicationModel.ServiceSession
	if req.GetSession().GetSessionId() != "" {
		sessionUser = applicationModel.SessionReq{
			SessionId:   uuid.MustParse(req.Session.GetSessionId()),
			UserId:      uuid.MustParse(req.Session.GetUserId()),
			CompanyId:   uuid.MustParse(req.Session.GetCompanyId()),
			ClientIp:    req.Session.GetClientIp(),
			ClientAgent: req.Session.GetClientAgent(),
		}
	}
	if req.GetServiceSession().GetServiceId() != "" {
		sessionService = applicationModel.ServiceSession{
			ServiceName: req.GetServiceSession().GetServiceName(),
			ServiceId:   req.GetServiceSession().GetServiceId(),
			ClientIp:    req.GetServiceSession().GetClientIp(),
			ClientAgent: req.GetServiceSession().GetClientAgent(),
		}
	}
	reqDelAttendanceEmployee := &applicationModel.DeleteAttendanceModel{
		Session:        &sessionUser,
		ServiceSession: &sessionService,
		//
		CompanyID: uuid.MustParse(req.GetCompanyId()),
		YearMonth: req.GetSummaryMonth(),
	}
	reqDelAttendanceRecordNoShift := &applicationModel.DeleteAttendanceRecordNoShiftModel{
		Session:        &sessionUser,
		ServiceSession: &sessionService,
		//
		CompanyID: uuid.MustParse(req.GetCompanyId()),
		YearMonth: req.GetSummaryMonth(),
	}
	// Handle del AttendanceRecords
	if err := s.attendanceService.DeleteAttendanceNoShift(
		ctx,
		reqDelAttendanceRecordNoShift,
	); err != nil {
		if err.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(500), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", err.ErrorClient)
	}
	// Handle del AttendanceEmployee
	if err := s.attendanceService.DeleteAttendanceRecord(
		ctx,
		reqDelAttendanceEmployee,
	); err != nil {
		if err.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(500), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", err.ErrorClient)
	}
	return &emptypb.Empty{}, nil
}

func (s *AttendanceGRPCServer) DeleteDailyAttendanceSummary(ctx context.Context, req *pb.DeleteAttendanceRecordsInput) (*emptypb.Empty, error) {
	// Parse request to application model
	var sessionUser applicationModel.SessionReq
	var sessionService applicationModel.ServiceSession
	if req.GetSession().GetSessionId() != "" {
		sessionUser = applicationModel.SessionReq{
			SessionId:   uuid.MustParse(req.Session.GetSessionId()),
			UserId:      uuid.MustParse(req.Session.GetUserId()),
			CompanyId:   uuid.MustParse(req.Session.GetCompanyId()),
			ClientIp:    req.Session.GetClientIp(),
			ClientAgent: req.Session.GetClientAgent(),
		}
	}
	if req.GetServiceSession().GetServiceId() != "" {
		sessionService = applicationModel.ServiceSession{
			ServiceName: req.GetServiceSession().GetServiceName(),
			ServiceId:   req.GetServiceSession().GetServiceId(),
			ClientIp:    req.GetServiceSession().GetClientIp(),
			ClientAgent: req.GetServiceSession().GetClientAgent(),
		}
	}
	repDeleteDailyAttendanceSummary := &applicationModel.DeleteDailyAttendanceSummaryModel{
		Session:        &sessionUser,
		ServiceSession: &sessionService,
		//
		CompanyID: uuid.MustParse(req.GetCompanyId()),
		SummaryMonth: req.GetSummaryMonth(),
	}
	// 
	if err := s.attendanceService.DeleteDailyAttendanceSummary(ctx,repDeleteDailyAttendanceSummary); err != nil {
		if err.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(500), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", err.ErrorClient)
	}
	return &emptypb.Empty{}, nil
}

func (s *AttendanceGRPCServer) HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *AttendanceGRPCServer) AddAttendance(ctx context.Context, req *pb.AddAttendanceInput) (*pb.AddAttendanceOutput, error) {
	// Mapping data request to service
	requestModel := &applicationModel.AddAttendanceModel{
		CompanyID:           uuid.MustParse(req.GetCompanyId()),
		EmployeeID:          uuid.MustParse(req.GetEmployeeId()),
		DeviceID:            uuid.MustParse(req.GetDeviceId()),
		RecordTime:          time.Unix(req.GetRecordTime(), 0),
		VerificationMethod:  req.GetVerificationMethod(),
		VerificationScore:   req.GetVerificationScore(),
		FaceImageURL:        req.GetFaceImageUrl(),
		LocationCoordinates: req.GetLocationCoordinates(),
		Session: &applicationModel.SessionReq{
			SessionId:   uuid.MustParse(req.Session.GetSessionId()),
			UserId:      uuid.MustParse(req.Session.GetUserId()),
			CompanyId:   uuid.MustParse(req.Session.GetCompanyId()),
			ClientIp:    req.Session.GetClientIp(),
			ClientAgent: req.Session.GetClientAgent(),
		},
	}
	// Call service
	appErr := s.attendanceService.AddAttendance(ctx, requestModel)
	if appErr != nil {
		if appErr.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(401), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", appErr.ErrorClient)
	}
	return &pb.AddAttendanceOutput{
		Message:    "Attendance added successfully",
		StatusCode: int32(codes.OK),
	}, nil
}

func (s *AttendanceGRPCServer) GetAttendanceRecords(ctx context.Context, input *pb.GetAttendanceRecordsInput) (*pb.GetAttendanceRecordsOutput, error) {
	// Mapping model req to application
	reqApplication := &applicationModel.GetAttendanceRecordsCompanyModel{
		CompanyID: uuid.MustParse(input.GetCompanyId()),
		YearMonth: input.GetYearMonth(),
		PageSize:  int(input.GetPageSize()),
		PageStage: []byte(input.GetPageStage()),
		Session: &applicationModel.SessionReq{
			SessionId:   uuid.MustParse(input.Session.GetSessionId()),
			UserId:      uuid.MustParse(input.Session.GetUserId()),
			CompanyId:   uuid.MustParse(input.Session.GetCompanyId()),
			ClientIp:    input.Session.GetClientIp(),
			ClientAgent: input.Session.GetClientAgent(),
		},
	}
	resp, err := s.attendanceService.GetAttendanceRecordsCompany(
		ctx,
		reqApplication,
	)
	if err != nil {
		if err.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(500), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", err.ErrorClient)
	}
	// Mapping response to grpc output
	output := &pb.GetAttendanceRecordsOutput{
		PageStageNext: []byte(resp.PageStageNext),
		PageSize:      int32(resp.PageSize),
		Records:       []*pb.AttendanceRecordInfo{},
	}
	for _, record := range resp.Records {
		output.Records = append(output.Records, &pb.AttendanceRecordInfo{
			CompanyId:           record.CompanyID.String(),
			YearMonth:           record.YearMonth,
			RecordTime:          record.RecordTime.Unix(),
			EmployeeId:          record.EmployeeID.String(),
			DeviceId:            record.DeviceID.String(),
			RecordType:          int32(record.RecordType),
			VerificationMethod:  record.VerificationMethod,
			VerificationScore:   record.VerificationScore,
			FaceImageUrl:        record.FaceImageURL,
			LocationCoordinates: record.LocationCoordinates,
			Metadata:            record.Metadata,
			SyncStatus:          record.SyncStatus,
			CreatedAt:           record.CreatedAt.Unix(),
		})
	}
	return output, nil
}

func (s *AttendanceGRPCServer) GetAttendanceRecordsEmployee(ctx context.Context, input *pb.GetAttendanceRecordsEmployeeInput) (*pb.GetAttendanceRecordsOutput, error) {
	// Mapping model req to application
	reqApplication := &applicationModel.GetAttendanceRecordsEmployeeModel{
		CompanyID:  uuid.MustParse(input.GetCompanyId()),
		YearMonth:  input.GetYearMonth(),
		EmployeeID: uuid.MustParse(input.GetEmployeeId()),
		PageSize:   int(input.GetPageSize()),
		PageStage:  []byte(input.GetPageStage()),
		Session: &applicationModel.SessionReq{
			SessionId:   uuid.MustParse(input.Session.GetSessionId()),
			UserId:      uuid.MustParse(input.Session.GetUserId()),
			CompanyId:   uuid.MustParse(input.Session.GetCompanyId()),
			ClientIp:    input.Session.GetClientIp(),
			ClientAgent: input.Session.GetClientAgent(),
		},
	}
	resp, err := s.attendanceService.GetAttendanceRecordsEmployeeForConpany(
		ctx,
		reqApplication,
	)
	if err != nil {
		if err.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(500), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", err.ErrorClient)
	}
	// Mapping response to grpc output
	output := &pb.GetAttendanceRecordsOutput{
		PageStageNext: []byte(resp.PageStageNext),
		PageSize:      int32(resp.PageSize),
		Records:       []*pb.AttendanceRecordInfo{},
	}
	for _, record := range resp.Records {
		output.Records = append(output.Records, &pb.AttendanceRecordInfo{
			CompanyId:           record.CompanyID.String(),
			YearMonth:           record.YearMonth,
			RecordTime:          record.RecordTime.Unix(),
			EmployeeId:          record.EmployeeID.String(),
			DeviceId:            record.DeviceID.String(),
			RecordType:          int32(record.RecordType),
			VerificationMethod:  record.VerificationMethod,
			VerificationScore:   record.VerificationScore,
			FaceImageUrl:        record.FaceImageURL,
			LocationCoordinates: record.LocationCoordinates,
			Metadata:            record.Metadata,
			SyncStatus:          record.SyncStatus,
			CreatedAt:           record.CreatedAt.Unix(),
		})
	}
	return output, nil
}

func (s *AttendanceGRPCServer) GetDailyAttendanceSummary(ctx context.Context, input *pb.GetDailyAttendanceSummaryInput) (*pb.GetDailyAttendanceSummaryOutput, error) {
	// Mapping model req to application
	reqApplication := &applicationModel.GetDailyAttendanceSummaryModel{
		Session: &applicationModel.SessionReq{
			SessionId:   uuid.MustParse(input.Session.GetSessionId()),
			UserId:      uuid.MustParse(input.Session.GetUserId()),
			CompanyId:   uuid.MustParse(input.Session.GetCompanyId()),
			ClientIp:    input.Session.GetClientIp(),
			ClientAgent: input.Session.GetClientAgent(),
		},
		CompanyID:    uuid.MustParse(input.GetCompanyId()),
		SummaryMonth: input.GetSummaryMonth(),
		WorkDate:     time.Unix(input.GetWorkDate(), 0),
		PageSize:     int(input.GetPageSize()),
		PageStage:    []byte(input.GetPageStage()),
	}
	resp, err := s.attendanceService.GetDailyAttendanceSummaryForCompany(
		ctx,
		reqApplication,
	)
	if err != nil {
		if err.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(500), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", err.ErrorClient)
	}
	// Mapping response to grpc output
	output := &pb.GetDailyAttendanceSummaryOutput{
		PageStageNext: []byte(resp.PageStageNext),
		PageSize:      int32(resp.PageSize),
		Records:       []*pb.DailyAttendanceSummaryInfo{},
	}
	for _, record := range resp.Records {
		output.Records = append(output.Records, &pb.DailyAttendanceSummaryInfo{
			CompanyId:         record.CompanyId.String(),
			SummaryMonth:      record.SummaryMonth,
			WorkDate:          record.WorkDate.Unix(),
			EmployeeId:        record.EmployeeId.String(),
			ShiftId:           record.ShiftId.String(),
			ActualCheckIn:     record.ActualCheckIn.Unix(),
			ActualCheckOut:    record.ActualCheckOut.Unix(),
			AttendanceStatus:  record.AttendanceStatus,
			LateMinutes:       int32(record.LateMinutes),
			EarlyLeaveMinutes: int32(record.EarlyLeaveMinutes),
			TotalWorkMinutes:  int32(record.TotalWorkMinutes),
			Notes:             record.Notes,
			UpdatedAt:         record.UpdatedAt.Unix(),
		})
	}
	return output, nil
}

func (s *AttendanceGRPCServer) GetDailyAttendanceSummaryEmployee(ctx context.Context, input *pb.GetDailyAttendanceSummaryEmployeeInput) (*pb.GetDailyAttendanceSummaryOutput, error) {
	// Mapping model req to application
	reqApplication := &applicationModel.GetDailyAttendanceSummaryEmployeeModel{
		Session: &applicationModel.SessionReq{
			SessionId:   uuid.MustParse(input.Session.GetSessionId()),
			UserId:      uuid.MustParse(input.Session.GetUserId()),
			CompanyId:   uuid.MustParse(input.Session.GetCompanyId()),
			ClientIp:    input.Session.GetClientIp(),
			ClientAgent: input.Session.GetClientAgent(),
		},
		CompanyID:    uuid.MustParse(input.GetCompanyId()),
		SummaryMonth: input.GetSummaryMonth(),
		PageSize:     int(input.GetPageSize()),
		PageStage:    []byte(input.GetPageStage()),
	}
	resp, err := s.attendanceService.GetDailyAttendanceSummaryEmployeeForCompany(
		ctx,
		reqApplication,
	)
	if err != nil {
		if err.ErrorSystem != nil {
			return nil, status.Errorf(codes.Code(500), "System is busy, please try again later")
		}
		return nil, status.Errorf(codes.Code(400), "%s", err.ErrorClient)
	}
	// Mapping response to grpc output
	output := &pb.GetDailyAttendanceSummaryOutput{
		PageStageNext: []byte(resp.PageStageNext),
		PageSize:      int32(resp.PageSize),
		Records:       []*pb.DailyAttendanceSummaryInfo{},
	}
	for _, record := range resp.Records {
		output.Records = append(output.Records, &pb.DailyAttendanceSummaryInfo{
			CompanyId:         record.CompanyId.String(),
			SummaryMonth:      record.SummaryMonth,
			WorkDate:          record.WorkDate.Unix(),
			EmployeeId:        record.EmployeeId.String(),
			ShiftId:           record.ShiftId.String(),
			ActualCheckIn:     record.ActualCheckIn.Unix(),
			ActualCheckOut:    record.ActualCheckOut.Unix(),
			AttendanceStatus:  record.AttendanceStatus,
			LateMinutes:       int32(record.LateMinutes),
			EarlyLeaveMinutes: int32(record.EarlyLeaveMinutes),
			TotalWorkMinutes:  int32(record.TotalWorkMinutes),
			Notes:             record.Notes,
			UpdatedAt:         record.UpdatedAt.Unix(),
		})
	}
	return output, nil
}

func (s *AttendanceGRPCServer) AddBatchAttendance(stream pb.AttendanceService_AddBatchAttendanceServer) error {
	ctx := stream.Context()
	sem := make(chan struct{}, 20)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 1)
	for {
		select {
		case err := <-errChan:
			return err
		default:
		}
		//
		req, err := stream.Recv()
		if err == io.EOF {
			break // Client đã gửi xong
		}
		if err != nil {
			return err
		}
		wg.Add(1)
		sem <- struct{}{} // Acquire token (block nếu đã đạt giới hạn 20)
		go func(req *pb.AddAttendanceInput) {
			defer wg.Done()
			defer func() { <-sem }()
			companyID, _ := uuid.Parse(req.GetCompanyId())
			employeeID, _ := uuid.Parse(req.GetEmployeeId())
			deviceID, _ := uuid.Parse(req.GetDeviceId())
			sessionID, _ := uuid.Parse(req.Session.GetSessionId())
			userID, _ := uuid.Parse(req.Session.GetUserId())
			sessionCompanyID, _ := uuid.Parse(req.Session.GetCompanyId())
			// Mapping data request to service
			requestModel := &applicationModel.AddAttendanceModel{
				CompanyID:           companyID,
				EmployeeID:          employeeID,
				DeviceID:            deviceID,
				RecordTime:          time.Unix(req.GetRecordTime(), 0),
				VerificationMethod:  req.GetVerificationMethod(),
				VerificationScore:   req.GetVerificationScore(),
				FaceImageURL:        req.GetFaceImageUrl(),
				LocationCoordinates: req.GetLocationCoordinates(),
				Session: &applicationModel.SessionReq{
					SessionId:   sessionID,
					UserId:      userID,
					CompanyId:   sessionCompanyID,
					ClientIp:    req.Session.GetClientIp(),
					ClientAgent: req.Session.GetClientAgent(),
				},
			}
			// Call service
			appErr := s.attendanceService.AddAttendance(ctx, requestModel)
			// Critical section: Gửi phản hồi hoặc lỗi
			mu.Lock()
			defer mu.Unlock()
			// Nếu đã có lỗi hệ thống xảy ra ở goroutine khác, không gửi thêm gì cả
			if len(errChan) > 0 {
				return
			}
			if appErr != nil {
				var grpcErr error
				if appErr.ErrorSystem != nil {
					grpcErr = status.Errorf(codes.Code(401), "System is busy, please try again later")
				} else {
					grpcErr = status.Errorf(codes.Code(400), "%s", appErr.ErrorClient)
				}
				// Đẩy lỗi vào channel để main loop biết và dừng lại
				select {
				case errChan <- grpcErr:
				default:
				}
				return
			}
			// Send response back to client
			if err := stream.SendMsg(&pb.AddAttendanceOutput{
				Message:    "Attendance added successfully",
				StatusCode: int32(codes.OK),
			}); err != nil {
				// Nếu gửi thất bại (ví dụ client ngắt kết nối)
				select {
				case errChan <- err:
				default:
				}
			}
		}(req)
	}
	// Đợi tất cả các request đang chạy hoàn tất
	wg.Wait()
	// Kiểm tra lần cuối xem có lỗi nào không
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (s *AttendanceGRPCServer) ServiceAddBatchAttendance(stream pb.AttendanceService_ServiceAddBatchAttendanceServer) error {
	ctx := stream.Context()
	sem := make(chan struct{}, 20)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 1)
	for {
		select {
		case err := <-errChan:
			return err
		default:
		}
		//
		req, err := stream.Recv()
		if err == io.EOF {
			break // Client đã gửi xong
		}
		if err != nil {
			return err
		}
		wg.Add(1)
		sem <- struct{}{} // Acquire token (block nếu đã đạt giới hạn 20)
		go func(req *pb.ServiceAddBatchAttendanceInput) {
			defer wg.Done()
			defer func() { <-sem }()
			companyID, _ := uuid.Parse(req.GetCompanyId())
			employeeID, _ := uuid.Parse(req.GetEmployeeId())
			deviceID, _ := uuid.Parse(req.GetDeviceId())
			// Mapping data request to service
			requestModel := &applicationModel.AddAttendanceModel{
				CompanyID:           companyID,
				EmployeeID:          employeeID,
				DeviceID:            deviceID,
				RecordTime:          time.Unix(req.GetRecordTime(), 0),
				VerificationMethod:  req.GetVerificationMethod(),
				VerificationScore:   req.GetVerificationScore(),
				FaceImageURL:        req.GetFaceImageUrl(),
				LocationCoordinates: req.GetLocationCoordinates(),
				ServiceSession: &applicationModel.ServiceSession{
					ServiceName: req.GetSession().GetServiceName(),
					ServiceId:   req.GetSession().GetServiceId(),
					ClientIp:    req.GetSession().GetClientIp(),
					ClientAgent: req.GetSession().GetClientAgent(),
				},
			}
			// Call service
			appErr := s.attendanceService.AddAttendance(ctx, requestModel)
			// Critical section: Gửi phản hồi hoặc lỗi
			mu.Lock()
			defer mu.Unlock()
			// Nếu đã có lỗi hệ thống xảy ra ở goroutine khác, không gửi thêm gì cả
			if len(errChan) > 0 {
				return
			}
			if appErr != nil {
				var grpcErr error
				if appErr.ErrorSystem != nil {
					grpcErr = status.Errorf(codes.Code(401), "System is busy, please try again later")
				} else {
					grpcErr = status.Errorf(codes.Code(400), "%s", appErr.ErrorClient)
				}
				// Đẩy lỗi vào channel để main loop biết và dừng lại
				select {
				case errChan <- grpcErr:
				default:
				}
				return
			}
			// Send response back to client
			if err := stream.SendMsg(&pb.AddAttendanceOutput{
				Message:    "Attendance added successfully",
				StatusCode: int32(codes.OK),
			}); err != nil {
				// Nếu gửi thất bại (ví dụ client ngắt kết nối)
				select {
				case errChan <- err:
				default:
				}
			}
		}(req)
	}
	// Đợi tất cả các request đang chạy hoàn tất
	wg.Wait()
	// Kiểm tra lần cuối xem có lỗi nào không
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
