package service

import (
	// "context"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/cache"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/logger"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/repository"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/cache"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/crypto"
)

const (
	// Cache key prefixes
	shiftCachePrefix = "shift:"
	shiftCacheTTL    = 3600 // 1 hour in seconds
)

// =================================================
// Shift service implementation interface
// =================================================
type ShiftService struct {
	shiftRepo        repository.IShiftRepository
	logger           logger.ILogger
	distributedCache cache.IDistributedCache
	localCache       cache.ILocalCache
}

// ChangeStatusShift implements service.IShiftService.
func (s *ShiftService) ChangeStatusShift(ctx context.Context, input *applicationModel.ChangeStatusShiftInput) *applicationError.Error {
	if input == nil {
		s.logger.Error("ChangeStatusShift - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("ChangeStatusShift - Start", "user_id", input.UserId, "shift_id", input.ShiftId, "is_active", input.IsActive)

	// Check permissions
	if input.Role == domainModel.RoleManager &&
		input.CompanyIdReq != input.CompanyId {
		s.logger.Error("ChangeStatusShift - Permission denied", "user_id", input.UserId, "company_id_req", input.CompanyIdReq, "company_id", input.CompanyId)
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("permission denied"),
			ErrorClient: "You do not have permission to change the status of this shift",
		}
	}
	var companyId uuid.UUID
	if input.Role == domainModel.RoleManager && input.CompanyIdReq == input.CompanyId {
		companyId = input.CompanyId
	} else {
		companyId = input.CompanyIdReq
	}

	// Call repository
	switch input.IsActive {
	case true:
		err := s.shiftRepo.EnableShiftWithId(
			ctx,
			&domainModel.EnableShiftInput{
				ShiftID:   input.ShiftId,
				CompanyId: companyId,
			},
		)
		if err != nil {
			s.logger.Error("ChangeStatusShift - Failed to activate shift", "error", err, "shift_id", input.ShiftId)
			return &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "Failed to activate shift",
			}
		}
	case false:
		err := s.shiftRepo.DisableShiftWithId(
			ctx,
			&domainModel.DisableShiftInput{
				ShiftID:   input.ShiftId,
				CompanyId: companyId,
			},
		)
		if err != nil {
			s.logger.Error("ChangeStatusShift - Failed to deactivate shift", "error", err, "shift_id", input.ShiftId)
			return &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "Failed to deactivate shift",
			}
		}
	}
	s.logger.Info("ChangeStatusShift - Success", "shift_id", input.ShiftId, "is_active", input.IsActive)

	return nil
}

// GetListShift implements service.IShiftService.
func (s *ShiftService) GetListShift(ctx context.Context, input *applicationModel.GetListShiftInput) ([]*applicationModel.GetDetailShiftOutput, *applicationError.Error) {
	if input == nil {
		s.logger.Error("GetListShift - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}
	s.logger.Info("GetListShift - Start", "user_id", input.UserId, "company_id", input.CompanyId)
	// Get data in cache
	key := utilsCache.GetKeyListShiftInCompany(
		utilsCrypto.GetHash(input.CompanyId.String()),
		input.Page,
	)
	cachedData := ""
	// Try to get from local cache first
	if data, err := s.localCache.Get(ctx, key); err == nil && data != "" {
		s.logger.Info("GetListShift - Cache hit (local)", "company_id", input.CompanyId, "page", input.Page)
		cachedData = data
	}
	// Try to get from distributed cache
	if cachedData == "" {
		if data, err := s.distributedCache.Get(ctx, key); err == nil && data != "" {
			s.logger.Info("GetListShift - Cache hit (distributed)", "company_id", input.CompanyId, "page", input.Page)
			cachedData = data
		}
	}
	if cachedData != "" {
		var output []*applicationModel.GetDetailShiftOutput
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &output); unmarshalErr == nil {
			return output, nil
		} else {
			s.logger.Warn("GetListShift - Failed to unmarshal cache data", "error", unmarshalErr)
		}
	}
	// Call repository
	limit := 20
	offset := (input.Page - 1) * limit
	reps, err := s.shiftRepo.ListShifts(
		ctx,
		&domainModel.ListShiftsInput{
			CompanyID: input.CompanyId,
			Limit:     int32(limit),
			Offset:    int32(offset),
		},
	)
	if err != nil {
		s.logger.Error("GetListShift - Failed to get shift list", "error", err, "company_id", input.CompanyId)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to get shift list",
		}
	}
	if len(reps) == 0 {
		s.logger.Info("GetListShift - No shifts found", "company_id", input.CompanyId)
		return []*applicationModel.GetDetailShiftOutput{}, nil
	}
	// Convert to application model
	var output []*applicationModel.GetDetailShiftOutput
	for _, shift := range reps {
		// Convert work days from []int32 to []int
		workDays := make([]int, len(shift.WorkDays))
		for i, day := range shift.WorkDays {
			workDays[i] = int(day)
		}
		output = append(output, &applicationModel.GetDetailShiftOutput{
			ShiftId:               shift.ShiftID.String(),
			CompanyId:             shift.CompanyID.String(),
			Name:                  shift.Name,
			Description:           shift.Description,
			StartTime:             shift.StartTime,
			EndTime:               shift.EndTime,
			BreakDurationMinutes:  int(shift.BreakDurationMinutes),
			GracePeriodMinutes:    int(shift.GracePeriodMinutes),
			EarlyDepartureMinutes: int(shift.EarlyDepartureMinutes),
			WorkDays:              workDays,
			IsActive:              shift.IsActive,
		})
	}

	// Cache the result
	if jsonData, jsonErr := json.Marshal(output); jsonErr == nil {
		// Store in distributed cache
		if setErr := s.distributedCache.SetTTL(ctx, key, string(jsonData), shiftCacheTTL); setErr != nil {
			s.logger.Warn("GetListShift - Failed to set distributed cache", "error", setErr)
		}
		// Store in local cache
		if setErr := s.localCache.SetTTL(ctx, key, string(jsonData), shiftCacheTTL); setErr != nil {
			s.logger.Warn("GetListShift - Failed to set local cache", "error", setErr)
		}
	}

	s.logger.Info("GetListShift - Success", "company_id", input.CompanyId, "page", input.Page)

	return output, nil
}

// CreateShift implements service.IShiftService.
func (s *ShiftService) CreateShift(ctx context.Context, input *applicationModel.CreateShiftInput) (*applicationModel.CreateShiftOutput, *applicationError.Error) {
	// Validate input
	if input == nil {
		s.logger.Error("CreateShift - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}
	var companyId uuid.UUID
	if input.Role == domainModel.RoleManager && input.CompanyIdReq == input.CompanyId {
		companyId = input.CompanyIdReq
	} else {
		companyId = input.CompanyIdReq
	}
	s.logger.Info("CreateShift - Start", "user_id", input.UserId, "company_id", companyId)

	// Convert work days from []int to []int32
	workDays := make([]int32, len(input.WorkDays))
	for i, day := range input.WorkDays {
		workDays[i] = int32(day)
	}

	// Create domain input
	domainInput := &domainModel.CreateShiftInput{
		CompanyID:             companyId,
		Name:                  input.Name,
		Description:           input.Description,
		StartTime:             input.StartTime,
		EndTime:               input.EndTime,
		BreakDurationMinutes:  int32(input.BreakDurationMinutes),
		GracePeriodMinutes:    int32(input.GracePeriodMinutes),
		EarlyDepartureMinutes: int32(input.EarlyDepartureMinutes),
		WorkDays:              workDays,
	}

	// Call repository
	shiftID, err := s.shiftRepo.CreateShift(ctx, domainInput)
	if err != nil {
		s.logger.Error("CreateShift - Failed to create shift", "error", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to create shift",
		}
	}

	s.logger.Info("CreateShift - Success", "shift_id", shiftID.String())

	return &applicationModel.CreateShiftOutput{
		ShiftId: shiftID.String(),
	}, nil
}

// DeleteShift implements service.IShiftService.
func (s *ShiftService) DeleteShift(ctx context.Context, input *applicationModel.DeleteShiftInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("DeleteShift - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("DeleteShift - Start", "user_id", input.UserId, "shift_id", input.ShiftId)

	// Call repository
	err := s.shiftRepo.DeleteShift(ctx, input.ShiftId)
	if err != nil {
		s.logger.Error("DeleteShift - Failed to delete shift", "error", err, "shift_id", input.ShiftId)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to delete shift",
		}
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("%s%s", shiftCachePrefix, input.ShiftId.String())

	// Delete from distributed cache
	if delErr := s.distributedCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("DeleteShift - Failed to delete from distributed cache", "error", delErr)
	}

	// Delete from local cache
	if delErr := s.localCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("DeleteShift - Failed to delete from local cache", "error", delErr)
	}

	s.logger.Info("DeleteShift - Success", "shift_id", input.ShiftId)

	return nil
}

// EditShift implements service.IShiftService.
func (s *ShiftService) EditShift(ctx context.Context, input *applicationModel.EditShiftInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("EditShift - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("EditShift - Start", "user_id", input.UserId, "shift_id", input.ShiftId)

	// Convert work days from []int to []int32
	workDays := make([]int32, len(input.WorkDays))
	for i, day := range input.WorkDays {
		workDays[i] = int32(day)
	}

	// Create domain input for update
	domainInput := &domainModel.UpdateTimeShiftInput{
		ShiftID:               input.ShiftId,
		StartTime:             input.StartTime,
		EndTime:               input.EndTime,
		BreakDurationMinutes:  int32(input.BreakDurationMinutes),
		GracePeriodMinutes:    int32(input.GracePeriodMinutes),
		EarlyDepartureMinutes: int32(input.EarlyDepartureMinutes),
		WorkDays:              workDays,
	}

	// Call repository
	err := s.shiftRepo.UpdateTimeShift(ctx, domainInput)
	if err != nil {
		s.logger.Error("EditShift - Failed to update shift", "error", err, "shift_id", input.ShiftId)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to update shift",
		}
	}

	// Invalidate cache after update
	cacheKey := fmt.Sprintf("%s%s", shiftCachePrefix, input.ShiftId.String())

	// Delete from distributed cache
	if delErr := s.distributedCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("EditShift - Failed to delete from distributed cache", "error", delErr)
	}

	// Delete from local cache
	if delErr := s.localCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("EditShift - Failed to delete from local cache", "error", delErr)
	}

	s.logger.Info("EditShift - Success", "shift_id", input.ShiftId)

	return nil
}

// GetDetailShift implements service.IShiftService.
func (s *ShiftService) GetDetailShift(ctx context.Context, input *applicationModel.GetDetailShiftInput) (*applicationModel.GetDetailShiftOutput, *applicationError.Error) {
	// Validate input
	if input == nil {
		s.logger.Error("GetDetailShift - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}
	s.logger.Info("GetDetailShift - Start", "user_id", input.UserId, "shift_id", input.ShiftId)

	cacheKey := fmt.Sprintf("%s%s", shiftCachePrefix, input.ShiftId.String())

	// Try to get from local cache first
	if cachedData, err := s.localCache.Get(ctx, cacheKey); err == nil && cachedData != "" {
		s.logger.Info("GetDetailShift - Cache hit (local)", "shift_id", input.ShiftId)
		var output applicationModel.GetDetailShiftOutput
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &output); unmarshalErr == nil {
			return &output, nil
		} else {
			s.logger.Warn("GetDetailShift - Failed to unmarshal local cache", "error", unmarshalErr)
		}
	}

	// Try to get from distributed cache
	if cachedData, err := s.distributedCache.Get(ctx, cacheKey); err == nil && cachedData != "" {
		s.logger.Info("GetDetailShift - Cache hit (distributed)", "shift_id", input.ShiftId)
		var output applicationModel.GetDetailShiftOutput
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &output); unmarshalErr == nil {
			// Store in local cache for faster access next time
			if jsonData, _ := json.Marshal(output); len(jsonData) > 0 {
				_ = s.localCache.SetTTL(ctx, cacheKey, string(jsonData), shiftCacheTTL)
			}
			return &output, nil
		} else {
			s.logger.Warn("GetDetailShift - Failed to unmarshal distributed cache", "error", unmarshalErr)
		}
	}

	s.logger.Info("GetDetailShift - Cache miss, fetching from database", "shift_id", input.ShiftId)

	// Call repository
	shift, err := s.shiftRepo.GetShiftByID(ctx, input.ShiftId)
	if err != nil {
		s.logger.Error("GetDetailShift - Failed to get shift details", "error", err, "shift_id", input.ShiftId)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to get shift details",
		}
	}

	// Convert work days from []int32 to []int
	workDays := make([]int, len(shift.WorkDays))
	for i, day := range shift.WorkDays {
		workDays[i] = int(day)
	}

	output := &applicationModel.GetDetailShiftOutput{
		ShiftId:               shift.ShiftID.String(),
		CompanyId:             shift.CompanyID.String(),
		Name:                  shift.Name,
		Description:           shift.Description,
		StartTime:             shift.StartTime,
		EndTime:               shift.EndTime,
		BreakDurationMinutes:  int(shift.BreakDurationMinutes),
		GracePeriodMinutes:    int(shift.GracePeriodMinutes),
		EarlyDepartureMinutes: int(shift.EarlyDepartureMinutes),
		WorkDays:              workDays,
		IsActive:              shift.IsActive,
	}

	// Cache the result
	if jsonData, jsonErr := json.Marshal(output); jsonErr == nil {
		// Store in distributed cache
		if setErr := s.distributedCache.SetTTL(ctx, cacheKey, string(jsonData), shiftCacheTTL); setErr != nil {
			s.logger.Warn("GetDetailShift - Failed to set distributed cache", "error", setErr)
		}
		// Store in local cache
		if setErr := s.localCache.SetTTL(ctx, cacheKey, string(jsonData), shiftCacheTTL); setErr != nil {
			s.logger.Warn("GetDetailShift - Failed to set local cache", "error", setErr)
		}
	}

	s.logger.Info("GetDetailShift - Success", "shift_id", input.ShiftId)

	return output, nil
}

// New instance
func NewShiftService() service.IShiftService {
	shiftRepo, err := repository.GetShiftRepository()
	if err != nil {
		panic(fmt.Sprintf("Failed to get shift repository: %v", err))
	}

	log := logger.GetLogger()
	if log == nil {
		panic("Failed to get logger instance")
	}

	distributedCache, err := cache.GetDistributedCache()
	if err != nil {
		panic(fmt.Sprintf("Failed to get distributed cache: %v", err))
	}

	localCache, err := cache.GetLocalCache()
	if err != nil {
		panic(fmt.Sprintf("Failed to get local cache: %v", err))
	}

	return &ShiftService{
		shiftRepo:        shiftRepo,
		logger:           log,
		distributedCache: distributedCache,
		localCache:       localCache,
	}
}
