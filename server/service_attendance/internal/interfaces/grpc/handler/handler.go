package handler

import (
	"context"
	"strconv"
	"time"

	pb "github.com/youknow2509/cio_verify_face/server/service_attendance/proto"

	"github.com/google/uuid"
	appModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	appService "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
	sharedUuid "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/uuid"
)

// ================================
// Grpc handlers
// ================================
type iHandler interface {
	// CheckIn handles check-in requests
	CheckIn(ctx context.Context, req *pb.CheckIOAttendanceRequest) (*pb.ResponseBase, error)
	// CheckOut handles check-out requests
	CheckOut(ctx context.Context, req *pb.CheckIOAttendanceRequest) (*pb.ResponseBase, error)
	// GetRecords retrieves attendance records
	GetRecords(ctx context.Context, req *pb.GetAttendanceRecordsRequest) (*pb.GetAttendanceRecordsResponse, error)
}

// handle struct implementing iHandler
type Handler struct{}

// CheckIn implements iHandler.
func (h *Handler) CheckIn(ctx context.Context, req *pb.CheckIOAttendanceRequest) (*pb.ResponseBase, error) {
	if req == nil || req.GetSessionInfo() == nil {
		return &pb.ResponseBase{Msg: "invalid request", StatusCode: 400}, nil
	}
	sess := req.GetSessionInfo()
	// parse uuids
	userUUID, err := sharedUuid.ParseUUID(sess.GetUserId())
	if err != nil {
		return &pb.ResponseBase{Msg: "invalid session user id", StatusCode: 400}, nil
	}
	sessionUUID, _ := sharedUuid.ParseUUID(sess.GetSessionId())
	deviceUUID, err := sharedUuid.ParseUUID(req.GetDeviceId())
	if err != nil {
		return &pb.ResponseBase{Msg: "invalid device id", StatusCode: 400}, nil
	}
	// recorded at - proto field is int64 unix seconds
	timestampStr := ""
	if req.GetRecordedAt() != 0 {
		timestampStr = strconv.FormatInt(req.GetRecordedAt(), 10)
	} else {
		// fallback to current time
		timestampStr = strconv.FormatInt(time.Now().Unix(), 10)
	}
	// build application input
	in := &appModel.CheckInInput{
		UserCheckInId:      uuid.Nil, // not provided in proto
		VerificationMethod: req.GetVerificationMethod(),
		VerificationScore:  req.GetVerificationScore(),
		FaceImageURL:       req.GetFaceImageUrl(),
		Timestamp:          timestampStr,
		Location:           req.GetLocation(),
		DeviceId:           deviceUUID,
		// session info
		UserID:      userUUID,
		SessionID:   sessionUUID,
		Role:        int(sess.GetRole()),
		ClientIp:    sess.GetClientIp(),
		ClientAgent: sess.GetClientAgent(),
	}
	if errRep := appService.GetAttendanceService().CheckInUser(ctx, in); errRep != nil {
		if errRep.ErrorSystem != nil {
			return &pb.ResponseBase{Msg: "internal server error", StatusCode: 500}, nil
		}
		return &pb.ResponseBase{Msg: errRep.ErrorClient, StatusCode: 400}, nil
	}
	return &pb.ResponseBase{Msg: "Check-in successful", StatusCode: 200}, nil
}

// CheckOut implements iHandler.
func (h *Handler) CheckOut(ctx context.Context, req *pb.CheckIOAttendanceRequest) (*pb.ResponseBase, error) {
	if req == nil || req.GetSessionInfo() == nil {
		return &pb.ResponseBase{Msg: "invalid request", StatusCode: 400}, nil
	}
	sess := req.GetSessionInfo()
	userUUID, err := sharedUuid.ParseUUID(sess.GetUserId())
	if err != nil {
		return &pb.ResponseBase{Msg: "invalid session user id", StatusCode: 400}, nil
	}
	sessionUUID, _ := sharedUuid.ParseUUID(sess.GetSessionId())
	deviceUUID, err := sharedUuid.ParseUUID(req.GetDeviceId())
	if err != nil {
		return &pb.ResponseBase{Msg: "invalid device id", StatusCode: 400}, nil
	}
	timestampStr := ""
	if req.GetRecordedAt() != 0 {
		timestampStr = strconv.FormatInt(req.GetRecordedAt(), 10)
	} else {
		timestampStr = strconv.FormatInt(time.Now().Unix(), 10)
	}
	in := &appModel.CheckOutInput{
		UserCheckOutId:     uuid.Nil,
		VerificationMethod: req.GetVerificationMethod(),
		VerificationScore:  req.GetVerificationScore(),
		FaceImageURL:       req.GetFaceImageUrl(),
		Timestamp:          timestampStr,
		Location:           req.GetLocation(),
		DeviceId:           deviceUUID,
		// session info
		UserID:      userUUID,
		SessionID:   sessionUUID,
		Role:        int(sess.GetRole()),
		ClientIp:    sess.GetClientIp(),
		ClientAgent: sess.GetClientAgent(),
	}
	if errRep := appService.GetAttendanceService().CheckOutUser(ctx, in); errRep != nil {
		if errRep.ErrorSystem != nil {
			return &pb.ResponseBase{Msg: "internal server error", StatusCode: 500}, nil
		}
		return &pb.ResponseBase{Msg: errRep.ErrorClient, StatusCode: 400}, nil
	}
	return &pb.ResponseBase{Msg: "Check-out successful", StatusCode: 200}, nil
}

// GetRecords implements iHandler.
func (h *Handler) GetRecords(ctx context.Context, req *pb.GetAttendanceRecordsRequest) (*pb.GetAttendanceRecordsResponse, error) {
	if req == nil || req.GetSessionInfo() == nil {
		return &pb.GetAttendanceRecordsResponse{Msg: "invalid request", StatusCode: 400}, nil
	}
	sess := req.GetSessionInfo()
	userUUID, err := sharedUuid.ParseUUID(sess.GetUserId())
	if err != nil {
		return &pb.GetAttendanceRecordsResponse{Msg: "invalid session user id", StatusCode: 400}, nil
	}
	sessionUUID, _ := sharedUuid.ParseUUID(sess.GetSessionId())
	// parse company and device
	var companyUUID uuid.UUID
	if req.GetCompanyId() != "" {
		companyUUID, err = sharedUuid.ParseUUID(req.GetCompanyId())
		if err != nil {
			return &pb.GetAttendanceRecordsResponse{Msg: "invalid company id", StatusCode: 400}, nil
		}
	} else {
		companyUUID = uuid.Nil
	}
	var deviceUUID uuid.UUID
	if req.GetDeviceId() != "" {
		deviceUUID, err = sharedUuid.ParseUUID(req.GetDeviceId())
		if err != nil {
			return &pb.GetAttendanceRecordsResponse{Msg: "invalid device id", StatusCode: 400}, nil
		}
	} else {
		deviceUUID = uuid.Nil
	}
	// convert date ranges
	var startDate time.Time
	var endDate time.Time
	if req.GetStartDate() != 0 {
		startDate = time.Unix(req.GetStartDate(), 0)
	}
	if req.GetEndDate() != 0 {
		endDate = time.Unix(req.GetEndDate(), 0)
	}
	// default pagination
	page := 1
	size := 20
	// call service
	recordsOut, errRep := appService.GetAttendanceService().GetRecords(
		ctx,
		&appModel.GetAttendanceRecordsInput{
			Page:      page,
			Size:      size,
			StartDate: startDate,
			EndDate:   endDate,
			DeviceID:  deviceUUID,
			CompanyId: companyUUID,
			// session
			UserID:      userUUID,
			SessionID:   sessionUUID,
			Role:        int(sess.GetRole()),
			ClientIp:    sess.GetClientIp(),
			ClientAgent: sess.GetClientAgent(),
		},
	)
	if errRep != nil {
		if errRep.ErrorSystem != nil {
			return &pb.GetAttendanceRecordsResponse{Msg: "internal server error", StatusCode: 500}, nil
		}
		return &pb.GetAttendanceRecordsResponse{Msg: errRep.ErrorClient, StatusCode: 400}, nil
	}
	// map service outputs to proto records
	var protoRecords []*pb.AttendanceRecord
	for _, out := range recordsOut {
		for _, info := range out.Records {
			// check-in
			if info.CheckIn != "" {
				if t, err := time.Parse(time.RFC3339, info.CheckIn); err == nil {
					protoRecords = append(protoRecords, &pb.AttendanceRecord{
						UserId:    info.UserID.String(),
						DeviceId:  info.DeviceID.String(),
						CheckType: "check_in",
						Timestamp: t.Unix(),
						Location:  info.Location,
					})
				} else {
					protoRecords = append(protoRecords, &pb.AttendanceRecord{
						UserId:    info.UserID.String(),
						DeviceId:  info.DeviceID.String(),
						CheckType: "check_in",
						Timestamp: 0,
						Location:  info.Location,
					})
				}
			}
			// check-out
			if info.CheckOut != "" {
				if t, err := time.Parse(time.RFC3339, info.CheckOut); err == nil {
					protoRecords = append(protoRecords, &pb.AttendanceRecord{
						UserId:    info.UserID.String(),
						DeviceId:  info.DeviceID.String(),
						CheckType: "check_out",
						Timestamp: t.Unix(),
						Location:  info.Location,
					})
				} else {
					protoRecords = append(protoRecords, &pb.AttendanceRecord{
						UserId:    info.UserID.String(),
						DeviceId:  info.DeviceID.String(),
						CheckType: "check_out",
						Timestamp: 0,
						Location:  info.Location,
					})
				}
			}
		}
	}
	return &pb.GetAttendanceRecordsResponse{Records: protoRecords, Msg: "ok", StatusCode: 200}, nil
}

// NewHandler creates a new instance of iHandler
func NewHandler() iHandler {
	return &Handler{}
}
