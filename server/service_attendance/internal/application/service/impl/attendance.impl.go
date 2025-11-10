package impl

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	applicationErrors "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/errors"
	model "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/cache"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/global"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/cache"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/crypto"
)

// ============================================
// Attendance Service
// ============================================
type AttendanceService struct{}

// CheckInUser implements service.IAttendanceService.
func (a *AttendanceService) CheckInUser(ctx context.Context, input *model.CheckInInput) *applicationErrors.Error {
	// Instance use
	var (
		repo             domainRepo.IAttendanceRepository
		userRepo         domainRepo.IUserRepository
		distributedCache domainCache.IDistributedCache
		logger           domainLogger.ILogger
	)
	logger = global.Logger
	distributedCache, _ = domainCache.GetDistributedCache()
	repo = domainRepo.GetAttendanceRepository()
	userRepo = domainRepo.GetUserRepository()
	// Check permission user manager or admin
	ok, err := checkUserManagerOrAdminForUser(
		ctx,
		input.UserID,
		input.CompanyId,
		input.UserCheckInId,
		input.Role,
	)
	if err != nil {
		return &applicationErrors.Error{ErrorSystem: err}
	}
	if !ok {
		return &applicationErrors.Error{ErrorClient: "User does not have permission to check in"}
	}
	// parse timestamp: try int unix seconds first, then RFC3339
	var recordTime int64
	if t, err := strconv.ParseInt(input.Timestamp, 10, 64); err == nil {
		recordTime = t
	} else if tm, err := time.Parse(time.RFC3339, input.Timestamp); err == nil {
		recordTime = tm.Unix()
	} else {
		return &applicationErrors.Error{ErrorClient: "Timestamp is not valid"}
	}

	// Prevent rapid duplicate check-ins for same user & date using distributed cache (atomic increment + expire)
	date := time.Unix(recordTime, 0).Format("2006-01-02")
	key := utilsCache.GetKeyAttendanceUserLastCheckIn(
		utilsCrypto.GetHash(input.UserID.String()),
		date,
	)
	// Try local cache first to avoid remote call
	if lc, err := domainCache.GetLocalCache(); err == nil {
		if ex, _ := lc.Exists(ctx, key); ex {
			return &applicationErrors.Error{ErrorClient: "Duplicate check-in detected"}
		}
	}
	// Atomic increment with TTL on distributed cache to avoid race conditions
	lua := `local v=redis.call('INCR', KEYS[1]); if tonumber(v)==1 then redis.call('EXPIRE', KEYS[1], tonumber(ARGV[1])) end; return v`
	res, err := distributedCache.LuaScript(ctx, lua, []string{key}, 60)
	if err == nil {
		// interpret result
		var val int64
		switch t := res.(type) {
		case int64:
			val = t
		case int:
			val = int64(t)
		case float64:
			val = int64(t)
		case string:
			if p, e := strconv.ParseInt(t, 10, 64); e == nil {
				val = p
			}
		}
		if val > 1 {
			return &applicationErrors.Error{ErrorClient: "Duplicate check-in detected"}
		}
		// set local cache to short TTL to prevent repeat checks on this instance
		if lc, err := domainCache.GetLocalCache(); err == nil {
			_ = lc.SetTTL(ctx, key, "1", 60)
		}
	} else {
		// fallback to Exists/SetTTL if LuaScript not supported
		if exists, e := distributedCache.Exists(ctx, key); e == nil && exists {
			return &applicationErrors.Error{ErrorClient: "Duplicate check-in detected"}
		}
		_ = distributedCache.SetTTL(ctx, key, "1", 60)
		if lc, err := domainCache.GetLocalCache(); err == nil {
			_ = lc.SetTTL(ctx, key, "1", 60)
		}
	}
	// get company ID of user
	companyReps, err := userRepo.GetCompanyIdUser(
		ctx,
		&domainModel.GetCompanyIdUserInput{
			UserID: input.UserCheckInId,
		},
	)
	if err != nil {
		logger.Error("GetCompanyIdUser failed", "error", err)
		return &applicationErrors.Error{ErrorSystem: err}
	}
	companyId := companyReps.CompanyID
	// build domain input and call repository
	repoInput := &domainModel.AddCheckInRecordInput{
		CompanyID:           companyId,
		EmployeeID:          input.UserCheckInId,
		DeviceID:            input.DeviceId,
		VerificationMethod:  input.VerificationMethod,
		VerificationScore:   input.VerificationScore,
		FaceImageURL:        input.FaceImageURL,
		LocationCoordinates: input.Location,
		Metadata: map[string]string{
			"client_ip":    input.ClientIp,
			"client_agent": input.ClientAgent,
		},
		RecordTime: recordTime,
	}
	if err := repo.AddCheckInRecord(ctx, repoInput); err != nil {
		logger.Error("AddCheckInRecord failed", "error", err)
		return &applicationErrors.Error{ErrorSystem: err}
	}
	logger.Info("Check-in recorded", "user_id", input.UserID.String(), "device_id", input.DeviceId.String())
	return nil
}

// CheckOutUser implements service.IAttendanceService.
func (a *AttendanceService) CheckOutUser(ctx context.Context, input *model.CheckOutInput) *applicationErrors.Error {
	// Instance use
	var (
		repo     domainRepo.IAttendanceRepository
		userRepo domainRepo.IUserRepository
		logger   domainLogger.ILogger
	)
	logger = global.Logger
	repo = domainRepo.GetAttendanceRepository()
	userRepo = domainRepo.GetUserRepository()
	// Check permission user manager or admin
	ok, err := checkUserManagerOrAdminForUser(
		ctx,
		input.UserID,
		input.CompanyId,
		input.UserCheckOutId,
		input.Role,
	)
	if err != nil {
		return &applicationErrors.Error{ErrorSystem: err}
	}
	if !ok {
		return &applicationErrors.Error{ErrorClient: "User does not have permission to check in"}
	}
	// parse timestamp
	var recordTime int64
	if t, err := strconv.ParseInt(input.Timestamp, 10, 64); err == nil {
		recordTime = t
	} else if tm, err := time.Parse(time.RFC3339, input.Timestamp); err == nil {
		recordTime = tm.Unix()
	} else {
		return &applicationErrors.Error{ErrorClient: "Timestamp is not valid"}
	}
	// get company ID of user
	companyReps, err := userRepo.GetCompanyIdUser(
		ctx,
		&domainModel.GetCompanyIdUserInput{
			UserID: input.UserCheckOutId,
		},
	)
	if err != nil {
		logger.Error("GetCompanyIdUser failed", "error", err)
		return &applicationErrors.Error{ErrorSystem: err}
	}
	companyId := companyReps.CompanyID
	// build domain input and call repository
	repoInput := &domainModel.AddCheckOutRecordInput{
		CompanyID:           companyId,
		EmployeeID:          input.UserCheckOutId,
		DeviceID:            input.DeviceId,
		VerificationMethod:  input.VerificationMethod,
		VerificationScore:   input.VerificationScore,
		FaceImageURL:        input.FaceImageURL,
		LocationCoordinates: input.Location,
		Metadata: map[string]string{
			"client_ip":    input.ClientIp,
			"client_agent": input.ClientAgent,
		},
		RecordTime: recordTime,
	}
	if err := repo.AddCheckOutRecord(ctx, repoInput); err != nil {
		logger.Error("AddCheckOutRecord failed", "error", err)
		return &applicationErrors.Error{ErrorSystem: err}
	}
	logger.Info("Check-out recorded", "user_id", input.UserID.String(), "device_id", input.DeviceId.String())
	return nil
}

// GetMyRecords implements service.IAttendanceService.
func (a *AttendanceService) GetMyRecords(ctx context.Context, input *model.GetMyRecordsInput) ([]*model.GetMyRecordsOutput, *applicationErrors.Error) {
	// Instance use
	var (
		repo             domainRepo.IAttendanceRepository
		distributedCache domainCache.IDistributedCache
		logger           domainLogger.ILogger
	)
	logger = global.Logger
	distributedCache, _ = domainCache.GetDistributedCache()
	repo = domainRepo.GetAttendanceRepository()
	// determine time range
	var start time.Time
	var end time.Time
	if input.StartDate.IsZero() {
		start = time.Now().Add(-30 * 24 * time.Hour)
	} else {
		start = input.StartDate
	}
	if input.EndDate.IsZero() {
		end = time.Now()
	} else {
		end = input.EndDate
	}
	// Try cache first (keyed by user + date range + page)
	cacheKey := ""
	if lc, err := domainCache.GetLocalCache(); err == nil {
		cacheKey = utilsCache.GetKeyAttendanceDeviceRecords(
			utilsCrypto.GetHash(input.CompanyId.String()),
			utilsCrypto.GetHash(input.UserID.String()),
			input.Page,
			input.Size,
			start.Format("2006-01-02"),
		)
		if v, e := lc.Get(ctx, cacheKey); e == nil && v != "" {
			var outCached []*model.GetMyRecordsOutput
			if jErr := json.Unmarshal([]byte(v), &outCached); jErr == nil {
				return outCached, nil
			}
		}
	}

	// build repo input; repository requires company and device - company not available here so use nil UUID
	repoInput := &domainModel.GetAttendanceRecordRangeTimeInput{
		CompanyID: input.CompanyId,
		DeviceID:  uuid.Nil,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
	}
	records, err := repo.GetAttendanceRecordRangeTime(ctx, repoInput)
	if err != nil {
		logger.Error("GetAttendanceRecordRangeTime failed", "error", err)
		return nil, &applicationErrors.Error{ErrorSystem: err}
	}

	// filter records only for current user and group by user into outputs
	grouped := map[uuid.UUID]*model.AttendanceRecordInfo{}
	for _, r := range records {
		if r.EmployeeID != input.UserID {
			continue
		}
		info, ok := grouped[r.EmployeeID]
		if !ok {
			info = &model.AttendanceRecordInfo{
				RecordID: uuid.Nil,
				UserID:   r.EmployeeID,
				DeviceID: r.DeviceID,
				Location: r.LocationCoordinates,
			}
			grouped[r.EmployeeID] = info
		}
		tStr := time.Unix(r.RecordTime, 0).Format(time.RFC3339)
		if r.Type == 0 {
			info.CheckIn = tStr
		} else {
			info.CheckOut = tStr
		}
	}
	out := &model.GetMyRecordsOutput{}
	for _, v := range grouped {
		out.Records = append(out.Records, *v)
	}
	out.Total = len(out.Records)
	result := []*model.GetMyRecordsOutput{out}
	// cache result in local and distributed caches
	if bs, jErr := json.Marshal(result); jErr == nil {
		if lc, err := domainCache.GetLocalCache(); err == nil {
			_ = lc.SetTTL(ctx, cacheKey, string(bs), 30)
		}
		if distributedCache != nil {
			_ = distributedCache.SetTTL(ctx, cacheKey, bs, 30)
		}
	}
	return result, nil
}

// GetRecords implements service.IAttendanceService.
func (a *AttendanceService) GetRecords(ctx context.Context, input *model.GetAttendanceRecordsInput) ([]*model.AttendanceRecordOutput, *applicationErrors.Error) {
	// Instance use
	var (
		repo             domainRepo.IAttendanceRepository
		distributedCache domainCache.IDistributedCache
		logger           domainLogger.ILogger
	)
	logger = global.Logger
	distributedCache, _ = domainCache.GetDistributedCache()
	repo = domainRepo.GetAttendanceRepository()
	// Check permission user manager or admin
	ok, err := checkUserManagerOrAdmin(
		ctx,
		input.UserID,
		input.CompanyIdUser,
		input.CompanyId,
		input.Role,
	)
	if err != nil {
		return nil, &applicationErrors.Error{ErrorSystem: err}
	}
	if !ok {
		return nil, &applicationErrors.Error{ErrorClient: "User does not have permission to check in"}
	}
	// determine time range
	var start time.Time
	var end time.Time
	if input.StartDate.IsZero() {
		start = time.Now().Add(-30 * 24 * time.Hour)
	} else {
		start = input.StartDate
	}
	if input.EndDate.IsZero() {
		end = time.Now()
	} else {
		end = input.EndDate
	}
	deviceId := input.DeviceID
	if deviceId == uuid.Nil {
		deviceId = uuid.Nil
	}
	// caching per device+date
	cacheKey := utilsCache.GetKeyAttendanceDeviceRecords(
		utilsCrypto.GetHash(input.CompanyId.String()),
		deviceId.String(),
		input.Page,
		input.Size,
		start.Format("2006-01-02"),
	)
	if lc, err := domainCache.GetLocalCache(); err == nil {
		if v, e := lc.Get(ctx, cacheKey); e == nil && v != "" {
			var outCached []*model.AttendanceRecordOutput
			if jErr := json.Unmarshal([]byte(v), &outCached); jErr == nil {
				return outCached, nil
			}
		}
	}

	repoInput := &domainModel.GetAttendanceRecordRangeTimeInput{
		CompanyID: input.CompanyId,
		DeviceID:  deviceId,
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		Limit:     input.Page,
		Offset:    (input.Page - 1) * input.Size,
	}
	records, err := repo.GetAttendanceRecordRangeTime(ctx, repoInput)
	if err != nil {
		logger.Error("GetAttendanceRecordRangeTime failed", "error", err)
		return nil, &applicationErrors.Error{ErrorSystem: err}
	}

	// group records by employee
	grouped := map[uuid.UUID]*model.AttendanceRecordInfo{}
	for _, r := range records {
		info, ok := grouped[r.EmployeeID]
		if !ok {
			info = &model.AttendanceRecordInfo{
				RecordID: uuid.Nil,
				UserID:   r.EmployeeID,
				DeviceID: r.DeviceID,
				Location: r.LocationCoordinates,
			}
			grouped[r.EmployeeID] = info
		}
		tStr := time.Unix(r.RecordTime, 0).Format(time.RFC3339)
		if r.Type == 0 {
			info.CheckIn = tStr
		} else {
			info.CheckOut = tStr
		}
	}
	out := &model.AttendanceRecordOutput{}
	for _, v := range grouped {
		out.Records = append(out.Records, *v)
	}
	out.Total = len(out.Records)
	result := []*model.AttendanceRecordOutput{out}
	if bs, jErr := json.Marshal(result); jErr == nil {
		if lc, err := domainCache.GetLocalCache(); err == nil {
			_ = lc.SetTTL(ctx, cacheKey, string(bs), 30)
		}
		if distributedCache != nil {
			_ = distributedCache.SetTTL(ctx, cacheKey, bs, 30)
		}
	}
	return result, nil
}

// New Attendance Service instance and impl interface
func NewAttendanceService() service.IAttendanceService {
	return &AttendanceService{}
}

// ============================================
//  Helper functions
// ============================================

// Check permission user manager or admin
func checkUserManagerOrAdmin(
	ctx context.Context,
	userReq uuid.UUID,
	companyIdUserReq uuid.UUID,
	companyId uuid.UUID,
	role int,
) (bool, error) {
	// Instance use
	var (
		localCache       domainCache.ILocalCache
		distributedCache domainCache.IDistributedCache
	)
	localCache, _ = domainCache.GetLocalCache()
	distributedCache, _ = domainCache.GetDistributedCache()
	// Check user role
	if role > domainModel.RoleManager {
		return false, nil
	}
	if role == domainModel.RoleAdmin {
		return true, nil
	}
	if companyIdUserReq != companyId {
		return false, nil
	}
	// Check cache first
	cacheKey := utilsCache.GetKeyUserIsManagerCompany(
		utilsCrypto.GetHash(userReq.String()),
		utilsCrypto.GetHash(companyId.String()),
	)
	if localCache != nil {
		if v, err := localCache.Get(ctx, cacheKey); err == nil && v == "1" {
			return true, nil
		}
	}
	if distributedCache != nil {
		if v, err := distributedCache.Get(ctx, cacheKey); err == nil && v == "1" {
			// set local cache for faster access next time
			if localCache != nil {
				_ = localCache.SetTTL(ctx, cacheKey, "1", 300)
			}
			return true, nil
		}
	}
	// Cache result if is manager
	if localCache != nil {
		_ = localCache.SetTTL(ctx, cacheKey, "1", 300)
	}
	if distributedCache != nil {
		_ = distributedCache.SetTTL(ctx, cacheKey, "1", 600)
	}
	return true, nil
}

// Check permission user manager or admin for user
func checkUserManagerOrAdminForUser(
	ctx context.Context,
	userReq uuid.UUID,
	userReqCompanyId uuid.UUID,
	user uuid.UUID,
	role int,
) (bool, error) {
	// Instance use
	var (
		repo             domainRepo.IUserRepository
		localCache       domainCache.ILocalCache
		distributedCache domainCache.IDistributedCache
	)
	localCache, _ = domainCache.GetLocalCache()
	distributedCache, _ = domainCache.GetDistributedCache()
	repo = domainRepo.GetUserRepository()
	// Check user role
	if role >= domainModel.RoleManager {
		return false, nil
	}
	if role == domainModel.RoleAdmin {
		return true, nil
	}
	// Check cache first
	cacheKey := utilsCache.GetKeyUserIsManagerCompanyForUser(
		utilsCrypto.GetHash(userReq.String()),
		utilsCrypto.GetHash(user.String()),
	)
	if localCache != nil {
		if v, err := localCache.Get(ctx, cacheKey); err == nil && v == "1" {
			return true, nil
		}
	}
	if distributedCache != nil {
		if v, err := distributedCache.Get(ctx, cacheKey); err == nil && v == "1" {
			// set local cache for faster access next time
			if localCache != nil {
				_ = localCache.SetTTL(ctx, cacheKey, "1", 300)
			}
			return true, nil
		}
	}
	// Check from repository
	isManager, err := repo.UserIsManagerCompany(
		ctx,
		&domainModel.UserIsManagerCompanyInput{
			UserID:    userReq,
			CompanyID: userReqCompanyId,
		},
	)
	if err != nil {
		return false, err
	}
	// Cache result if is manager
	if isManager {
		if localCache != nil {
			_ = localCache.SetTTL(ctx, cacheKey, "1", 300)
		}
		if distributedCache != nil {
			_ = distributedCache.SetTTL(ctx, cacheKey, "1", 600)
		}
	}
	return isManager, nil
}
